package ql

import (
	"fmt"
	"strconv"

	"github.com/SethCurry/stax/internal/bones/card"
	"github.com/SethCurry/stax/internal/bones/cardface"
	"github.com/SethCurry/stax/internal/bones/predicate"
)

// FieldFilterHandler is a function that can return an AST node for
// a given value.  These are typically registered with a FieldFilter.
type FieldFilterHandler func(value string) (leaf, error)

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

func FloatFieldFilterHandler(handler func(value float32) (leaf, error)) FieldFilterHandler {
	return func(value string) (leaf, error) {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float: %w", err)
		}

		return handler(float32(f))
	}
}

// colorsNotInQuery returns a slice of colors that are not in the query.
// This is useful for constructing some of the comparison operators for colors
// and color identities.
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

var allColors = []string{"W", "U", "B", "R", "G"}

type colorQuery struct {
	MustInclude []string
	MustExclude []string
}

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

// FieldFilterHandlerWithColors is a helper function that converts a string of colors
// into a slice of colors and calls the given handler.
// You can use this to create FieldFilterHandlers with the colors already parsed.
func FieldFilterHandlerWithColors(handler func([]string) (leaf, error)) FieldFilterHandler {
	return func(colorString string) (leaf, error) {
		var colors []string

		for _, value := range colorString {
			colors = append(colors, string(value))
		}

		return handler(colors)
	}
}

func colorsEQ() FieldFilterHandler {
	return FieldFilterHandlerWithColors(func(colors []string) (leaf, error) {
		colorQuery := &colorQuery{
			MustInclude: colors,
			MustExclude: colorsNotInQuery(colors),
		}
		return &basicLeaf{predicator: colorQuery.toPredicator()}, nil
	})
}

func colorsLT() FieldFilterHandler {
	return FieldFilterHandlerWithColors(func(colors []string) (leaf, error) {
		colorQuery := &colorQuery{
			MustExclude: colorsNotInQuery(colors),
		}
		return &basicLeaf{predicator: colorQuery.toPredicator()}, nil
	})
}

func colorsGT() FieldFilterHandler {
	return FieldFilterHandlerWithColors(func(colors []string) (leaf, error) {
		colorQuery := &colorQuery{
			MustInclude: colors,
		}
		return &basicLeaf{predicator: colorQuery.toPredicator()}, nil
	})
}
