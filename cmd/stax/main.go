package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SethCurry/stax/internal/cli"
	"github.com/alecthomas/kong"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	userHome, err := homedir.Dir()
	if err != nil {
		fmt.Println("Failed to get your home directory.")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var root cli.Root

	configPaths := []string{
		filepath.Join(userHome, ".stax", "config.json"),
	}

	ctx := kong.Parse(&root,
		kong.Name("stax"),
		kong.Description("A swiss army knife CLI for Magic: The Gathering"),
		kong.Configuration(kong.JSON, configPaths...))

	var logLevel zapcore.Level
	levelRef := &logLevel

	err = levelRef.UnmarshalText([]byte(root.Verbosity))
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
