package responses

import (
	"github.com/SethCurry/stax/internal/bones"
	"github.com/SethCurry/stax/internal/common"
)

type Card struct {
	Name     string     `json:"name"`
	OracleID string     `json:"oracle_id"`
	Faces    []CardFace `json:"faces"`
}

type CardFace struct {
	Name       string  `json:"name"`
	OracleText string  `json:"oracle_text"`
	FlavorText string  `json:"flavor_text"`
	Language   string  `json:"language"`
	CMC        float32 `json:"cmc"`
	Power      string  `json:"power"`
	Toughness  string  `json:"toughness"`
	Loyalty    string  `json:"loyalty"`
	ManaCost   string  `json:"mana_cost"`
	TypeLine   string  `json:"type_line"`
	Colors     string  `json:"colors"`
}

type CardSearch struct {
	Cards []Card `json:"cards"`
}

// CardFromDB converts a single Card to a Card response object.
func CardFromDB(crd *bones.Card) Card {
	return Card{
		Name:     crd.Name,
		OracleID: crd.OracleID,
		Faces: common.Map(crd.Edges.Faces, func(f *bones.CardFace) CardFace {
			return CardFace{
				Name:       f.Name,
				OracleText: f.OracleText,
				FlavorText: f.FlavorText,
				Language:   f.Language,
				CMC:        f.Cmc,
				Power:      f.Power,
				Toughness:  f.Toughness,
				Loyalty:    f.Loyalty,
				ManaCost:   f.ManaCost,
				TypeLine:   f.TypeLine,
				Colors:     f.Colors,
			}
		}),
	}
}

// CardsFromDB converts a slice of database cards to Card response objects.
func CardsFromDB(crds []*bones.Card) []Card {
	return common.Map(crds, CardFromDB)
}
