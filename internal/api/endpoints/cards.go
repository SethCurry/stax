package endpoints

import (
	"errors"
	"fmt"

	"github.com/SethCurry/stax/internal/api/squid"
	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"github.com/SethCurry/stax/internal/oracle/oracledb/card"
	"github.com/SethCurry/stax/pkg/common"
)

type CardByNameQuery struct {
	Exact string `schema:"exact"`
	Fuzzy string `schema:"fuzzy"`
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

	result, err := query.WithFaces().Only(ctx.Request.Context())
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
