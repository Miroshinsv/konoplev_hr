package main

import (
	"context"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"

	"github.com/meBazil/hr-mvp/internal"
	"github.com/meBazil/hr-mvp/internal/config"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
	"github.com/meBazil/hr-mvp/pkg/process-manager"
)

var (
	cfg config.App
)

func init() {
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zlog := zerolog.New(os.Stdout).Level(level)
	defaultLogger := logger.NewLoggerFrom(context.Background(), &zlog)
	logger.SetDefaultLogger(defaultLogger)

	process.SetLogger(&PMLogger{Logger: &defaultLogger})
}

func main() {
	application, err := internal.NewApplication(cfg)
	if err != nil {
		panic(err)
	}

	application.Run()
}
