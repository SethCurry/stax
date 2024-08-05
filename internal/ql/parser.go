package ql

import (
	"fmt"

	"github.com/SethCurry/stax/internal/bones/card"
	"github.com/SethCurry/stax/internal/bones/predicate"
)

type Operator string

const (
	EQ Operator = "="
	NE Operator = "!="
	GT Operator = ">"
	GE Operator = ">="
	LT Operator = "<"
	LE Operator = "<="
)

type filterType func(Operator, string) (Leaf, error)

type cardFilter func(string) Leaf

type Leaf interface {
	Predicate() predicate.Card
}

type Node interface {
	Leaf
	Left() Leaf
	Right() Leaf
}

type LogicNode struct {
	predicator func(...predicate.Card) predicate.Card
	left       Leaf
	right      Leaf
}

func (l *LogicNode) Predicate() predicate.Card {
	return l.predicator(l.left.Predicate(), l.right.Predicate())
}

func (l *LogicNode) Left() Leaf {
	return l.left
}

func (l *LogicNode) Right() Leaf {
	return l.right
}

type basicLeaf struct {
	predicator predicate.Card
}

func (l *basicLeaf) Predicate() predicate.Card {
	return l.predicator
}

func cardNameSearch(op Operator, name string) (Leaf, error) {
	switch op {
	case EQ:
		return &basicLeaf{
			predicator: card.NameEQ(name),
		}, nil
	}

	return nil, fmt.Errorf("invalid operator: %s", op)
}

var filterTypeLookupTable = map[string]filterType{
	"name": cardNameSearch,
}
