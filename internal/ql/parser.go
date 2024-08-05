package ql

import (
	"errors"
	"fmt"

	"github.com/SethCurry/scurry-go/fp"
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
	SetLeft(Leaf)
	SetRight(Leaf)
}

type LogicNode struct {
	predicator func(...predicate.Card) predicate.Card
	left       Leaf
	right      Leaf
}

func (l *LogicNode) Predicate() predicate.Card {
	var leftValue predicate.Card
	var rightValue predicate.Card

	if l.left != nil {
		leftValue = l.left.Predicate()
	}

	if l.right != nil {
		rightValue = l.right.Predicate()
	}

	return l.predicator(leftValue, rightValue)
}

func (l *LogicNode) Left() Leaf {
	return l.left
}

func (l *LogicNode) Right() Leaf {
	return l.right
}

func (l *LogicNode) SetLeft(left Leaf) {
	l.left = left
}

func (l *LogicNode) SetRight(right Leaf) {
	l.right = right
}

func newAndNode(left, right Leaf) *LogicNode {
	return &LogicNode{
		predicator: func(cards ...predicate.Card) predicate.Card {
			return card.And(fp.Filter[predicate.Card](func(c predicate.Card) bool {
				return c != nil
			}, cards)...)
		},
		left:  left,
		right: right,
	}
}

func newOrNode(left, right Leaf) *LogicNode {
	return &LogicNode{
		predicator: func(cards ...predicate.Card) predicate.Card {
			return card.Or(cards...)
		},
		left:  left,
		right: right,
	}
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
			predicator: card.NameContainsFold(name),
		}, nil
	}

	return nil, fmt.Errorf("invalid operator: %s", op)
}

var filterTypeLookupTable = map[string]filterType{
	"name": cardNameSearch,
}

func getLeaf(filterType string, op Operator, value string) (Leaf, error) {
	filterTypeFunc, ok := filterTypeLookupTable[filterType]
	if !ok {
		return nil, fmt.Errorf("invalid filter type: %s", filterType)
	}

	return filterTypeFunc(op, value)
}

func newTokenReader(tokens []Token) *tokenReader {
	return &tokenReader{
		tokens: tokens,
		index:  0,
	}
}

type tokenReader struct {
	tokens []Token
	index  int
}

func (r *tokenReader) next() (*Token, bool) {
	// return nil if the index is out of bounds
	// i.e. there are no more tokens to read
	if !r.hasMore() {
		return nil, false
	}

	ret := r.tokens[r.index]
	r.index++

	return &ret, true
}

func (r *tokenReader) peek() (*Token, bool) {
	if !r.hasMore() {
		return nil, false
	}

	return &r.tokens[r.index], true
}

func (r *tokenReader) hasMore() bool {
	return r.index < len(r.tokens)
}

func parseTokens(tokens []Token) (Node, error) {
	reader := newTokenReader(tokens)

	var root Node
	var previous Node

	for reader.hasMore() {
		nextToken, ok := reader.next()
		if !ok {
			return nil, errors.New("unexpected EOF")
		}

		switch nextToken.Family {
		case FamilyOperator:
			return nil, fmt.Errorf("expected keyword instead of operator, but got: %s", nextToken.Value)
		case FamilyParen:
			return nil, errors.New("parentheses are not yet supported")
		case FamilyKeyword:
			return nil, errors.New("keywords are not yet supported")
		case FamilyLiteral:
			filterType := nextToken.Value

			opToken, ok := reader.next()
			if !ok {
				return nil, errors.New("expected operator after literal")
			}

			valueToken, ok := reader.next()
			if !ok {
				return nil, errors.New("expected value after operator")
			}

			leafNode, err := getLeaf(filterType, Operator(opToken.Value), valueToken.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to build leaf node: %w", err)
			}

			if root == nil {
				root = newAndNode(leafNode, previous)
				previous = root
			} else {
				newPrevious := newAndNode(nil, nil)
				previous.SetLeft(leafNode)
				previous.SetRight(newPrevious)
				previous = newPrevious
			}
		}
	}

	return root, nil
}

func ParseQuery(query string) (Node, error) {
	tokens, err := Lex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to lex query: %w", err)
	}

	return parseTokens(tokens)
}
