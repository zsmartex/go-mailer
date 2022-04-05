package main

import (
	"github.com/zsmartex/go-mailer/internal/config"
	"github.com/zsmartex/go-mailer/pkg/consumer"
)

func main() {
	config.InitConfig()

	consumer := consumer.NewConsumer()

	consumer.Run()
}
