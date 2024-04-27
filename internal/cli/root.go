package cli

import "go.uber.org/zap"

type Context struct {
	Logger *zap.Logger
}

// Root is the main entrypoint for the CLI.
// All of the subcommands start here.
type Root struct {
	// The scryfall command.
	Scryfall ScryfallCmd `cmd:"" help:"Scryfall API commands"`

	// The log level to use.
	// This needs to be unmarshaled into a zapcore.Level.
	LogLevel string `name:"log-level" help:"The level to log at." default:"error" enum:"debug,info,warn,error"`
}
