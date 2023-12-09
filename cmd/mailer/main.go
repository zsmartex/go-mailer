package main

import (
	"github.com/zsmartex/pkg/v2/infrastructure/context_fx"
	"go.uber.org/fx"

	"github.com/zsmartex/go-mailer/internal/config"
	"github.com/zsmartex/go-mailer/pkg/consumer"
)

func main() {
	app := fx.New(
		config.Module,
		context_fx.Module,
		consumer.Module,
	)

	app.Run()
}
