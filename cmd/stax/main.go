package main

import (
	"context"
	"fmt"

	"github.com/SethCurry/stax/internal/cli"
	"github.com/alecthomas/kong"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var root cli.Root

	ctx := kong.Parse(&root)

	var logLevel zapcore.Level
	levelRef := &logLevel

	err := levelRef.UnmarshalText([]byte(root.LogLevel))
	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	loggerConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Encoding:         "console",
		Development:      true,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "lvl",
			TimeKey:     "ts",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			EncodeTime:  zapcore.RFC3339TimeEncoder,
		},
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		ctx.FatalIfErrorf(err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	err = ctx.Run(&cli.Context{
		Logger:  logger,
		Context: context.Background(),
	})
	if err != nil {
		logger.Fatal("failed to execute command", zap.Error(err))
	}
}
