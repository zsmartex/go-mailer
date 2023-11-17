package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/cockroachdb/errors"
	"github.com/gookit/goutil/fsutil"
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"

	"github.com/zsmartex/go-mailer/pkg/eventapi"
)

var Module = fx.Module("config_fx", fx.Provide(NewConfig, provideConfig))

// Config represents application configuration model.
type Config struct {
	Kafka    config.Kafka
	Keychain map[string]eventapi.Validator `yaml:"keychain"`
	Topics   map[string]Topic              `yaml:"topics"`
	Events   []Event                       `yaml:"events"`
}

func NewConfig() (*Config, error) {
	var cfgYaml *struct {
		Keychain map[string]eventapi.Validator `yaml:"keychain"`
		Topics   map[string]Topic              `yaml:"topics"`
		Events   []Event                       `yaml:"events"`
	}
	mailerBytes := fsutil.MustReadFile("config/mailer.yml")

	err := yaml.Unmarshal(mailerBytes, &cfgYaml)
	if err != nil {
		return nil, err
	}

	conf := new(Config)

	if err := env.Parse(conf); err != nil {
		return nil, errors.Newf("parse config: %v", err)
	}

	conf.Keychain = cfgYaml.Keychain
	conf.Topics = cfgYaml.Topics
	conf.Events = cfgYaml.Events

	return conf, nil
}

type OutConfig struct {
	fx.Out

	Kafka config.Kafka
}

func provideConfig(config *Config) OutConfig {
	return OutConfig{
		Kafka: config.Kafka,
	}
}
