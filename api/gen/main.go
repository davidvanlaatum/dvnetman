package main

import (
	"context"
	"dvnetman/api/gen/code"
	"dvnetman/api/gen/spec"
	"dvnetman/pkg/logger"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	log := logger.NewLogger(logger.LevelTrace, logger.NewStdoutDriver(logger.NewConsoleFormatter().EnableColor(true)))
	ctx = log.Context(ctx)
	api := spec.NewSpec()
	if err := api.WriteOpenAPISpec(ctx, "api/openapi.yaml"); err != nil {
		panic(err)
	}
	gen := code.NewCodeGen()
	if err := gen.Generate(ctx, &api.OpenAPI); err != nil {
		panic(err)
	}
	if err := gen.WriteFiles(ctx, "pkg/openapi"); err != nil {
		panic(err)
	}
}
