package consumer

import (
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

type Email struct {
	FromAddress string
	FromName    string
	ToAddress   string
	Subject     string
	Content     string
}

type SMTPConf struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (conf SMTPConf) URL() string {
	return fmt.Sprintf("%s:%s", conf.Host, conf.Port)
}

type EmailSender struct {
	Conf  *SMTPConf
	Email *Email
}

func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Name: str, Address: ""}
	return strings.Trim(addr.String(), " <>")
}

func NewEmailSender(conf SMTPConf, email Email) *EmailSender {
	return &EmailSender{&conf, &email}
}

func (e *EmailSender) Send() error {
	// Password is required.
	if strings.TrimSpace(e.Conf.Password) == "" {
		return errors.New("password is empty")
	}

	if e.Email == nil {
		return errors.New("email is nil")
	}

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf(`%s <%s>`, e.Email.FromName, e.Email.FromAddress)
	headers["To"] = e.Email.ToAddress
	headers["Subject"] = encodeRFC2047(e.Email.Subject)
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + string(e.Email.Content)

	recipients := []string{e.Email.ToAddress}

	auth := smtp.PlainAuth("", e.Conf.Username, e.Conf.Password, e.Conf.Host)
	if err := smtp.SendMail(e.Conf.URL(), auth, e.Email.FromAddress, recipients, []byte(message)); err != nil {
		return err
	}

	return nil
}
