package requests

import (
	"errors"

	"github.com/SethCurry/stax/internal/bones/card"
	"github.com/SethCurry/stax/internal/bones/predicate"
)

type CardByName struct {
	Exact string `schema:"exact"`
	Fuzzy string `schema:"fuzzy"`
}

func (c CardByName) Validate() error {
	if c.Exact != "" && c.Fuzzy != "" {
		return errors.New("exact and fuzzy cannot be used at the same time")
	}

	if c.Exact == "" && c.Fuzzy == "" {
		return errors.New("either fuzzy or exact must be specified")
	}

	return nil
}

func (c CardByName) ToPredicate() predicate.Card {
	if c.Exact != "" {
		return card.NameEQ(c.Exact)
	}

	if c.Fuzzy != "" {
		return card.NameContainsFold(c.Fuzzy)
	}

	return card.And()
}

type CardSearch struct {
	Name string `schema:"name"`
}

type CardQuery struct {
	Query string `schema:"q"`
}
