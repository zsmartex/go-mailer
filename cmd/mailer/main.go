package main

import (
	"github.com/zsmartex/pkg/v2/infrastructure/context_fx"
	"github.com/zsmartex/pkg/v2/infrastructure/kafka_fx"
	"go.uber.org/fx"

	"github.com/zsmartex/go-mailer/internal/config"
	"github.com/zsmartex/go-mailer/pkg/consumer"
)

func main() {
	app := fx.New(
		config.Module,
		context_fx.Module,
		kafka_fx.ConsumerModule,
		consumer.Module,
	)

	app.Run()
}
