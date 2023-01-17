package consumer

import (
	"bytes"
	"context"
	"html"
	"os"
	"strings"
	"text/template"

	"github.com/gookit/goutil/fsutil"
	"gopkg.in/yaml.v3"

	"github.com/zsmartex/go-mailer/internal/config"
	"github.com/zsmartex/go-mailer/pkg/eventapi"
	"github.com/zsmartex/pkg/v2/infrastucture/kafka"
	"github.com/zsmartex/pkg/v2/log"
)

type Consumer struct {
	config *config.Config
}

func NewConsumer() *Consumer {
	var config *config.Config
	mailerBytes := fsutil.MustReadFile("config/mailer.yml")

	yaml.Unmarshal(mailerBytes, &config)

	return &Consumer{
		config: config,
	}
}

func (c *Consumer) Run() {
	topics := make([]string, 0)
	for _, exchange := range c.config.Topics {
		topics = append(topics, exchange.Name)
	}

	ctx := context.Background()

	log.Info("Starting mailer...")

	borkers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	consumer, err := kafka.NewConsumer(ctx, borkers, "zsmartex-mailer", topics...)
	if err != nil {
		log.Panicf("Failed to connect to kafka brokers err: %v", err)
	}

	for {
		records, err := consumer.Poll(ctx)
		if err != nil {
			log.Fatalf("Failed to poll from consumer err: %v", err)
		}

		for _, record := range records {
			log.Debugf("Received event: %s", string(record.Key))
			var eventConf *config.Event
			for _, e := range c.config.Events {
				if bytes.Equal(record.Key, []byte(e.Key)) {
					eventConf = &e
					break
				}
			}

			if eventConf == nil {
				log.Warnf("Not found event for key: %s", string(record.Key))
				consumer.CommitRecords(ctx, *record)
				continue
			}

			signer := c.config.Topics[eventConf.Topic].Signer

			c.handleEvent(eventConf, string(record.Value), signer)
			consumer.CommitRecords(ctx, *record)
		}
	}
}

func (c *Consumer) handleEvent(eventConf *config.Event, payload, signer string) error {
	validator := c.config.Keychain[signer]
	claims, err := eventapi.ParseJWT(payload, validator.ValidateJWT)
	if err != nil {
		log.Errorf("Failed to parse jwt err: %v", err)
		return err
	}

	event, err := eventapi.Unmarshal(claims.Event)
	if err != nil {
		log.Errorf("Failed to unmarshal event err: %v", err)
		return err
	}

	record, err := event.FixAndValidate(os.Getenv("DEFAULT_LANGUAGE"))
	if err != nil {
		log.Errorf("Failed to validate event err: %v", err)
		return err
	}

	tpl := eventConf.Template(record.Language)

	claims.Event["logo"] = os.Getenv("SENDER_LOGO") // set logo url
	body, err := tpl.Content(claims.Event)
	if err != nil {
		log.Errorf("Failed to execution template err: %v", err)
		return err
	}

	layout_tpl := template.Must(template.ParseFiles("templates/layout.tpl"))
	var buf bytes.Buffer
	if err := layout_tpl.Execute(&buf, map[string]interface{}{"Body": body}); err != nil {
		log.Errorf("Failed to execute template err: %v", err)
		return err
	}

	content := buf.String()
	content = html.UnescapeString(content)

	email := Email{
		FromAddress: os.Getenv("SENDER_EMAIL"),
		FromName:    os.Getenv("SENDER_NAME"),
		ToAddress:   record.User.Email,
		Subject:     tpl.Subject,
		Content:     string(content),
	}

	password := os.Getenv("SMTP_PASSWORD")
	conf := SMTPConf{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USER"),
		Password: password,
	}

	if err := NewEmailSender(conf, email).Send(); err != nil {
		log.Errorf("Failed to send email: %v", err)

		return err
	}

	log.Debugf("Sent email to: %s", email.ToAddress)

	return nil
}
