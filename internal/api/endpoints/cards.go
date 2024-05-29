// Package endpoints contains the HTTP handlers for the stax API.
package endpoints

import (
	"errors"
	"fmt"

	"github.com/SethCurry/stax/internal/api/squid"
	"github.com/SethCurry/stax/internal/common"
	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"github.com/SethCurry/stax/internal/oracle/oracledb/card"
	"github.com/SethCurry/stax/internal/oracle/oracledb/predicate"
)

type CardByNameQuery struct {
	Exact string `schema:"exact"`
	Fuzzy string `schema:"fuzzy"`
}

func (c CardByNameQuery) Validate() error {
	if c.Exact != "" && c.Fuzzy != "" {
		return errors.New("exact and fuzzy cannot be used at the same time")
	}

	if c.Exact == "" && c.Fuzzy == "" {
		return errors.New("either fuzzy or exact must be specified")
	}

	return nil
}

func (c CardByNameQuery) ToPredicate() predicate.Card {
	if c.Exact != "" {
		return card.NameEQ(c.Exact)
	}

	if c.Fuzzy != "" {
		return card.NameContainsFold(c.Fuzzy)
	}

	return card.And()
}

func cardToResponse(crd *oracledb.Card) CardResponse {
	return CardResponse{
		Name:     crd.Name,
		OracleID: crd.OracleID,
		Faces: common.Map(crd.Edges.Faces, func(f *oracledb.CardFace) CardFace {
			return CardFace{
				OracleText: f.OracleText,
			}
		}),
	}
}

func cardsToResponse(crds []*oracledb.Card) []CardResponse {
	return common.Map(crds, cardToResponse)
}

func CardByName(ctx *squid.Context) error {
	var params CardByNameQuery

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

	resp := cardToResponse(result)

	return ctx.Response.WriteJSON(200, resp)
}

type CardSearchRequest struct {
	Name string `schema:"name"`
}

type CardFace struct {
	OracleText string `json:"oracle_text"`
}

type CardResponse struct {
	Name     string     `json:"name"`
	OracleID string     `json:"oracle_id"`
	Faces    []CardFace `json:"faces"`
}

type CardSearchResponse struct {
	Cards []CardResponse `json:"cards"`
}

func CardSearch(ctx *squid.Context) error {
	var params CardSearchRequest

	if err := ctx.Request.UnmarshalQuery(&params); err != nil {
		return err
	}

	query := ctx.DB.Card.Query()

	if params.Name != "" {
		query = query.Where(card.NameContainsFold(params.Name))
	}

	results, err := query.WithFaces().All(ctx.Request.Context())
	if err != nil {
		return fmt.Errorf("failed to query for cards: %w", err)
	}

	cards := cardsToResponse(results)

	return ctx.Response.WriteJSON(200, CardSearchResponse{
		Cards: cards,
	})
}
