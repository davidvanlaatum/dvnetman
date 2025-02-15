package main

import (
	"dvnetman/api/gen/code"
	"dvnetman/api/gen/spec"
	"dvnetman/pkg/logger"
)

func main() {
	log := logger.NewLogger(logger.LevelTrace, logger.NewStdoutDriver(logger.NewConsoleFormatter().EnableColor(true)))
	api := spec.NewSpec(log)
	if err := api.WriteOpenAPISpec("api/openapi.yaml"); err != nil {
		panic(err)
	}
	gen := code.NewCodeGen(log)
	if err := gen.Generate(&api.OpenAPI); err != nil {
		panic(err)
	}
	if err := gen.WriteFiles("pkg/openapi"); err != nil {
		panic(err)
	}
}
