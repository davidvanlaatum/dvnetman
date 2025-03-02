package main

import (
	"context"
	"dvnetman/pkg/config"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/server"
	"errors"
	"flag"
	"os"
	"os/signal"
)

func main() {
	configPath := flag.String("config", "/etc/dvnetman.yaml", "config file")

	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	log := logger.NewLogger(logger.LevelTrace, logger.NewStdoutDriver(logger.NewConsoleFormatter().EnableColor(true)))
	ctx = log.SubLogger().KeyProvider(logger.OTelTraceKeyProvider).Logger().Context(ctx)

	if *configPath == "" {
		logger.Error(ctx).Msg("config file is required")
		os.Exit(1)
	}

	if cfg, err := config.LoadConfig(*configPath); err != nil {
		logger.Error(ctx).Err(err).Msg("failed to load config")
	} else {
		var s = server.NewServer(cfg)
		if err = s.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Fatal(ctx).Err(err).Msg("failed to start server")
		}
	}
}
