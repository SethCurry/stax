package endpoints

import (
	"errors"
	"fmt"

	"github.com/SethCurry/stax/internal/api/squid"
	"github.com/SethCurry/stax/internal/oracle/oracledb/card"
)

type CardByNameQuery struct {
	Exact string `schema:"exact"`
	Fuzzy string `schema:"fuzzy"`
}

type CardResponse struct {
	Name string `json:"name"`
}

func CardByName(ctx *squid.Context) error {
	var params CardByNameQuery

	if err := ctx.Request.UnmarshalQuery(&params); err != nil {
		return err
	}

	if params.Exact == "" && params.Fuzzy == "" {
		return errors.New("must provide either fuzzy or exact")
	}

	if params.Exact != "" && params.Fuzzy != "" {
		return errors.New("can provide either fuzzy or exact, not both")
	}

	query := ctx.DB.Card.Query()

	if params.Exact != "" {
		query = query.Where(card.NameEQ(params.Exact))
	}

	if params.Fuzzy != "" {
		query = query.Where(card.NameContainsFold(params.Fuzzy))
	}

	result, err := query.Only(ctx.Request.Context())
	if err != nil {
		return fmt.Errorf("failed to query card: %w", err)
	}

	resp := CardResponse{
		Name: result.Name,
	}

	return ctx.Response.WriteJSON(200, resp)
}
