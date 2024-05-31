package responses

import (
	"github.com/SethCurry/stax/internal/common"
	"github.com/SethCurry/stax/internal/oracle/oracledb"
)

type Card struct {
	Name     string     `json:"name"`
	OracleID string     `json:"oracle_id"`
	Faces    []CardFace `json:"faces"`
}

type CardFace struct {
	OracleText string `json:"oracle_text"`
}

type CardSearch struct {
	Cards []Card `json:"cards"`
}

// CardFromDB converts a single Card to a Card response object.
func CardFromDB(crd *oracledb.Card) Card {
	return Card{
		Name:     crd.Name,
		OracleID: crd.OracleID,
		Faces: common.Map(crd.Edges.Faces, func(f *oracledb.CardFace) CardFace {
			return CardFace{
				OracleText: f.OracleText,
			}
		}),
	}
}

// CardsFromDB converts a slice of database cards to Card response objects.
func CardsFromDB(crds []*oracledb.Card) []Card {
	return common.Map(crds, CardFromDB)
}
