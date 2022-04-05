package config

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zsmartex/go-mailer/pkg/eventapi"
	"github.com/zsmartex/pkg/services"
)

// Template represents email massage content and subject.
type Template struct {
	Subject      string `yaml:"subject"`
	TemplatePath string `yaml:"template_path,omitempty"`
	Template     string `yaml:"template,omitempty"`
}

// Event represent configuration for listening an message from RabbitMQ.
type Event struct {
	Name       string              `yaml:"name"`
	Key        string              `yaml:"key"`
	Topic      string              `yaml:"topic"`
	Templates  map[string]Template `yaml:"templates"`
	Expression string              `yaml:"expression"`
}

// Exchange contains exchange name and signer unique identifier.
type Topic struct {
	Name   string `yaml:"name"`
	Signer string `yaml:"signer"`
}

// Config represents application configuration model.
type Config struct {
	Keychain map[string]eventapi.Validator `yaml:"keychain"`
	Topics   map[string]Topic              `yaml:"topics"`
	Events   []Event                       `yaml:"events"`
}

func (e *Event) Template(key string) Template {
	return e.Templates[key]
}

// Content returns ready to go message with specified data.
// Note: "template" has bigger priority, than "template_path".
func (t *Template) Content(data interface{}) (string, error) {
	var err error

	buff := new(bytes.Buffer)
	var tpl *template.Template

	if strings.TrimSpace(t.Template) != "" {
		tpl, err = template.New(t.Subject).Parse(t.Template)
	} else {
		tpl, err = template.ParseFiles(t.TemplatePath)
	}

	if err != nil {
		return "", err
	}

	if err := tpl.Execute(buff, &data); err != nil {
		return "", err
	}

	return buff.String(), nil
}

var Logger *logrus.Entry

func InitConfig() (err error) {
	Logger = services.NewLoggerService("Mailer")

	return nil
}
