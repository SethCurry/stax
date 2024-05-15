package cli

import (
	"context"

	"go.uber.org/zap"
)

type Context struct {
	Logger  *zap.Logger
	Context context.Context
}

// Root is the main entrypoint for the CLI.
// All of the subcommands start here.
type Root struct {
	// The scryfall command.
	Scryfall ScryfallCmd `cmd:"" help:"Scryfall API commands"`

	Oracle OracleCmd `cmd:"" help:"Oracle database commands"`

	API APICmd `cmd:"" help:"Start the API server."`

	// The log level to use.
	// This needs to be unmarshaled into a zapcore.Level.
	LogLevel string `name:"log-level" help:"The level to log at." default:"error" enum:"debug,info,warn,error"`
}
