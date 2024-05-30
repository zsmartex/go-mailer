package consumer

import (
	"bytes"
	"context"
	"html"
	"os"
	"text/template"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/zsmartex/pkg/v2/infrastructure/kafka_fx"
	"github.com/zsmartex/pkg/v2/log"
	"go.uber.org/fx"

	"github.com/zsmartex/go-mailer/internal/config"
	"github.com/zsmartex/go-mailer/pkg/eventapi"
)

var (
	Module = fx.Module("consumer_fx",
		fx.Supply(kafka_fx.Group("mailer")),
		fx.Supply(time.NewTicker(time.Second)),
		fx.Supply(fx.Annotate(true, fx.ParamTags(`name:"at_end"`))),
		kafka_fx.ConsumerModule,
		fx.Provide(fx.Annotate(NewConsumer, fx.As(new(kafka_fx.ConsumerSubscriber)))),
		fx.Invoke(registerHooks),
	)
)

var _ kafka_fx.ConsumerSubscriber = (*Consumer)(nil)

type Consumer struct {
	config *config.Config
}

func NewConsumer(config *config.Config) *Consumer {
	return &Consumer{
		config: config,
	}
}

func (c *Consumer) OnMessage(record *kgo.Record) error {
	log.Debugf("Received event: %s", string(record.Key))
	var eventConf *config.Event
	for _, e := range c.config.Events {
		if bytes.Equal(record.Key, []byte(e.Key)) {
			eventConf = &e
			break
		}
	}

	if eventConf == nil {
		return nil
	}

	signer := c.config.Topics[eventConf.Topic].Signer

	return c.handleEvent(eventConf, string(record.Value), signer)
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

	var toAddress string
	emailData, ok := record.Data["email"]
	if ok {
		toAddress = emailData.(string)
	} else {
		toAddress = record.User.Email
	}

	email := Email{
		FromAddress: os.Getenv("SENDER_EMAIL"),
		FromName:    os.Getenv("SENDER_NAME"),
		ToAddress:   toAddress,
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

func registerHooks(config *config.Config, kafkaConsumer *kafka_fx.Consumer) error {
	topics := make([]kafka_fx.Topic, 0)
	for _, exchange := range config.Topics {
		topics = append(topics, kafka_fx.Topic(exchange.Name))
	}

	return kafkaConsumer.AddConsumeTopics(context.Background(), topics...)
}
