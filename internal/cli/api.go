package cli

import (
	"github.com/SethCurry/stax/internal/api/squid"
	"go.uber.org/zap"
)

type ApiCmd struct{}

func (a *ApiCmd) Run(ctx *Context) error {
	srv := squid.NewServer(ctx.Logger)

	err := srv.Serve(":8765")
	if err != nil {
		ctx.Logger.Error("failed to serve HTTP", zap.Error(err))
	}

	return nil
}
