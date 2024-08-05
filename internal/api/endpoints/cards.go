// Package endpoints contains the HTTP handlers for the stax API.
package endpoints

import (
	"fmt"

	"github.com/SethCurry/stax/internal/api/requests"
	"github.com/SethCurry/stax/internal/api/responses"
	"github.com/SethCurry/stax/internal/api/squid"
	"github.com/SethCurry/stax/internal/ql"
)

// CardByName searched for a single card via its name.  It should
// only ever return a single result.
func CardByName(ctx *squid.Context) error {
	var params requests.CardByName

	if err := ctx.Request.UnmarshalQuery(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	result, err := ctx.DB.Card.Query().Where(params.ToPredicate()).WithFaces().Only(ctx.Request.Context())
	if err != nil {
		return fmt.Errorf("failed to query card: %w", err)
	}

	resp := responses.CardFromDB(result)

	return ctx.Response.WriteJSON(200, resp)
}

func CardSearch(ctx *squid.Context) error {
	var params requests.CardQuery

	if err := ctx.Request.UnmarshalQuery(&params); err != nil {
		return err
	}

	parsedQueryRoot, err := ql.ParseQuery(params.Query)
	if err != nil {
		return err
	}

	pred := parsedQueryRoot.Predicate()

	gotCards, err := ctx.DB.Card.Query().Where(pred).WithFaces().All(ctx.Request.Context())
	if err != nil {
		return fmt.Errorf("failed to query for cards: %w", err)
	}

	cards := responses.CardsFromDB(gotCards)

	return ctx.Response.WriteJSON(200, responses.CardSearch{
		Cards: cards,
	})
}
