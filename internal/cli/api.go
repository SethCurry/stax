package cli

import (
	"context"

	"github.com/SethCurry/stax/internal/api/endpoints"
	"github.com/SethCurry/stax/internal/api/squid"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type APICmd struct {
	Listen string `optional:"" aliases:"l" name:"listen" default:"0.0.0.0:8765" help:"The address to listen on."`
}

func (a *APICmd) Run(ctx *Context) error {
	dbConn, err := connectToDatabase(context.Background(), ctx.Logger, false)
	if err != nil {
		ctx.Logger.Fatal("failed to open connection to DB", zap.Error(err))
	}

	srv := squid.NewServer(dbConn, ctx.Logger)

	srv.Get("/cards/named", endpoints.CardByName)
	srv.Get("/cards", endpoints.CardSearch)

	err = srv.Serve(a.Listen)
	if err != nil {
		ctx.Logger.Error("failed to serve HTTP", zap.Error(err))
	}

	return nil
}
