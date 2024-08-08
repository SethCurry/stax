package ql

import (
	"fmt"

	"github.com/SethCurry/stax/internal/bones/card"
	"github.com/SethCurry/stax/internal/bones/cardface"
	"github.com/SethCurry/stax/internal/bones/predicate"
)

type FieldFilterHandler func(value string) (leaf, error)

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

type FieldFilter struct {
	Name     string
	Aliases  []string
	Handlers map[operator]FieldFilterHandler
}

func (f *FieldFilter) Register(op operator, handler FieldFilterHandler) {
	f.Handlers[op] = handler
}

func (f *FieldFilter) MatchesName(name string) bool {
	if f.Name == name {
		return true
	}

	for _, alias := range f.Aliases {
		if alias == name {
			return true
		}
	}

	return false
}

func (f *FieldFilter) Handle(op operator, value string) (leaf, error) {
	if handler, ok := f.Handlers[op]; ok {
		return handler(value)
	}

	return nil, &ErrNoOperationForField{
		Field:    f.Name,
		Operator: op,
	}
}

// ErrNoOperationForField is returned when a field does not support a given operator.
// This occurs when parsing a query like "name < Drannith", because there is no
// less than operator for the name field.
type ErrNoOperationForField struct {
	Field    string
	Operator operator
}

func (e *ErrNoOperationForField) Error() string {
	return fmt.Sprintf("field %q has no such operation: %s", e.Field, e.Operator)
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
