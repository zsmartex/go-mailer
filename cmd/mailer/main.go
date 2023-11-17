package main

import (
	"go.uber.org/fx"

	"github.com/zsmartex/go-mailer/internal/config"
	"github.com/zsmartex/go-mailer/pkg/consumer"
)

func main() {
	app := fx.New(
		config.Module,
		consumer.Module,
	)

	app.Run()
}
