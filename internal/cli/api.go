package cli

import (
	"github.com/SethCurry/stax/internal/api/squid"
	"github.com/SethCurry/stax/internal/oracle/oracledb"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type ApiCmd struct{}

func (a *ApiCmd) Run(ctx *Context) error {
	dbConn, err := oracledb.Open("sqlite3", "file://.local/api.sqlite3")
	if err != nil {
		ctx.Logger.Fatal("failed to open connection to DB", zap.Error(err))
	}

	srv := squid.NewServer(dbConn, ctx.Logger)

	err = srv.Serve(":8765")
	if err != nil {
		ctx.Logger.Error("failed to serve HTTP", zap.Error(err))
	}

	return nil
}
