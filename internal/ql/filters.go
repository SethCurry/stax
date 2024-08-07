package ql

import (
	"fmt"
	"strconv"

	"github.com/SethCurry/stax/internal/bones/card"
	"github.com/SethCurry/stax/internal/bones/cardface"
	"github.com/SethCurry/stax/internal/bones/predicate"
)

func colorsNotInQuery(queryColors []string) []string {
	notInQuery := []string{}

	for _, color := range allColors {
		found := false

		for _, queryColor := range queryColors {
			if color == queryColor {
				found = true
				break
			}
		}

		if !found {
			notInQuery = append(notInQuery, color)
		}
	}

	return notInQuery
}

type colorQuery struct {
	MustInclude []string
	MustExclude []string
}

var allColors = []string{"W", "U", "B", "R", "G"}

func (c *colorQuery) toPredicator() predicate.Card {
	preds := []predicate.Card{}

	for _, include := range c.MustInclude {
		preds = append(preds, card.HasFacesWith(cardface.ColorsContainsFold(include)))
	}

	for _, exclude := range c.MustExclude {
		preds = append(preds, card.HasFacesWith(cardface.Not(cardface.ColorsContainsFold(exclude))))
	}

	return card.And(preds...)
}

func colorsSearch(op operator, value string) (leaf, error) {
	queryColors := []string{}

	for _, c := range value {
		queryColors = append(queryColors, string(c))
	}

	switch op {
	case opEQ:
		colorQuery := &colorQuery{
			MustInclude: queryColors,
			MustExclude: colorsNotInQuery(queryColors),
		}
		return &basicLeaf{
			predicator: colorQuery.toPredicator(),
		}, nil
	case opLT, opLE:
		colorQuery := &colorQuery{
			MustExclude: colorsNotInQuery(queryColors),
		}
		return &basicLeaf{
			predicator: colorQuery.toPredicator(),
		}, nil
	case opGT, opGE:
		colorQuery := &colorQuery{
			MustInclude: queryColors,
		}
		return &basicLeaf{
			predicator: colorQuery.toPredicator(),
		}, nil
	}

	return nil, fmt.Errorf("invalid operator: %s", op)
}

func cardNameSearch(op operator, name string) (leaf, error) {
	switch op {
	case opEQ:
		return &basicLeaf{
			predicator: card.NameContainsFold(name),
		}, nil
	}

	return nil, fmt.Errorf("invalid operator: %s", op)
}

func cmcSearch(op operator, value string) (leaf, error) {
	parsedValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %s", value)
	}

	switch op {
	case opEQ:
		return &basicLeaf{
			predicator: card.HasFacesWith(cardface.CmcEQ(float32(parsedValue))),
		}, nil
	case opLT:
		return &basicLeaf{
			predicator: card.HasFacesWith(cardface.CmcLT(float32(parsedValue))),
		}, nil
	case opLE:
		return &basicLeaf{
			predicator: card.HasFacesWith(cardface.CmcLTE(float32(parsedValue))),
		}, nil
	case opGT:
		return &basicLeaf{
			predicator: card.HasFacesWith(cardface.CmcGT(float32(parsedValue))),
		}, nil
	case opGE:
		return &basicLeaf{
			predicator: card.HasFacesWith(cardface.CmcGTE(float32(parsedValue))),
		}, nil
	}

	return nil, fmt.Errorf("invalid operator: %s", op)
}

func oracleTextSearch(op operator, value string) (leaf, error) {
	switch op {
	case opEQ:
		return &basicLeaf{
			predicator: card.HasFacesWith(cardface.OracleTextContainsFold(value)),
		}, nil
	}

	return nil, fmt.Errorf("invalid operator: %s", op)
}

var filterFieldLookupTable = map[string]filterField{
	"name":   cardNameSearch,
	"o":      oracleTextSearch,
	"oracle": oracleTextSearch,
	"text":   oracleTextSearch,
	"cmc":    cmcSearch,
	"colors": colorsSearch,
	"c":      colorsSearch,
}

func getLeaf(filterType string, op operator, value string) (leaf, error) {
	filterTypeFunc, ok := filterFieldLookupTable[filterType]
	if !ok {
		return nil, fmt.Errorf("invalid filter type: %s", filterType)
	}

	return filterTypeFunc(op, value)
}
