// Package app represents the application logic.
package app

import (
	"github.com/venafi/vmware-avi-connector/internal/app/discovery"
	vmwareavi "github.com/venafi/vmware-avi-connector/internal/app/vmware-avi"
	"github.com/venafi/vmware-avi-connector/internal/handler/web"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New ...
func New() *fx.App {
	var logger *zap.Logger

	app := fx.New(
		fx.Provide(
			configureLogger,
			web.ConfigureHTTPServers,
			fx.Annotate(vmwareavi.NewVMwareAviClients, fx.As(new(vmwareavi.ClientServices))),
			fx.Annotate(discovery.NewDiscoveryService, fx.As(new(vmwareavi.DiscoveryService))),
			fx.Annotate(vmwareavi.NewWebhookService, fx.As(new(web.WebhookService))),
		),
		fx.Invoke(
			web.RegisterHandlers,
		),
		fx.Populate(&logger),
	)

	logger.Info("VMware AVI connector starting")

	return app
}

func configureLogger() (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	loggerConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	loggerConfig.EncoderConfig.TimeKey = "time"
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(zap.L())
	return zap.L(), nil
}
