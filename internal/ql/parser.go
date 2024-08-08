package ql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/SethCurry/scurry-go/fp"
	"github.com/SethCurry/stax/internal/bones/card"
	"github.com/SethCurry/stax/internal/bones/cardface"
	"github.com/SethCurry/stax/internal/bones/predicate"
)

type operator string

const (
	opEQ operator = "="
	opNE operator = "!="
	opGT operator = ">"
	opGE operator = ">="
	opLT operator = "<"
	opLE operator = "<="
)

type filterField func(operator, string) (leaf, error)

type leaf interface {
	Predicate() predicate.Card
}

type node interface {
	leaf
	Left() leaf
	Right() leaf
	SetLeft(leaf)
	SetRight(leaf)
}

type logicNode struct {
	predicator func(...predicate.Card) predicate.Card
	left       leaf
	right      leaf
}

func (l *logicNode) Predicate() predicate.Card {
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

func (l *logicNode) Left() leaf {
	return l.left
}

func (l *logicNode) Right() leaf {
	return l.right
}

func (l *logicNode) SetLeft(left leaf) {
	l.left = left
}

func (l *logicNode) SetRight(right leaf) {
	l.right = right
}

func newAndNode(left, right leaf) *logicNode {
	return &logicNode{
		predicator: func(cards ...predicate.Card) predicate.Card {
			return card.And(fp.Filter[predicate.Card](func(c predicate.Card) bool {
				return c != nil
			}, cards)...)
		},
		left:  left,
		right: right,
	}
}

func newOrNode(left, right leaf) *logicNode {
	return &logicNode{
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

func (r *tokenReader) hasMore() bool {
	return r.index < len(r.tokens)
}

type Parser struct {
	Fields []FieldFilter
}

func (p *Parser) AddField(field FieldFilter) {
	p.Fields = append(p.Fields, field)
}

type ErrNoField struct {
	Field string
}

func (e *ErrNoField) Error() string {
	return fmt.Sprintf("no such field: %s", e.Field)
}

func (p *Parser) handleField(field string, op operator, value string) (leaf, error) {
	for _, f := range p.Fields {
		if f.MatchesName(field) {
			return f.Handle(op, value)
		}
	}

	return nil, &ErrNoField{Field: field}
}

// ParseTokens parses a slice of tokens and returns a node that can be converted to a bones predicate.
// This is useful if you want to separate the lexing and parsing phases.
func (p *Parser) ParseTokens(tokens []Token) (node, error) {
	reader := newTokenReader(tokens)

	var root node
	var previous node

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
			keyword := nextToken.Value

			switch strings.ToLower(keyword) {
			case "and":
				newPrevious := newAndNode(nil, nil)
				previous.SetRight(newPrevious)
				previous = newPrevious
			case "or":
				newPrevious := newOrNode(nil, nil)
				previous.SetRight(newPrevious)
				previous = newPrevious
			default:
				return nil, fmt.Errorf("unrecognized keyword: %s", keyword)
			}
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

			leafNode, err := p.handleField(filterType, operator(opToken.Value), valueToken.Value)
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

func (p *Parser) ParseQuery(query string) (node, error) {
	tokens, err := LexString(query)
	if err != nil {
		return nil, fmt.Errorf("failed to lex query: %w", err)
	}

	return p.ParseTokens(tokens)
}

func floatFieldFilterHandler(handler func(value float32) (leaf, error)) FieldFilterHandler {
	return func(value string) (leaf, error) {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float: %w", err)
		}

		return handler(float32(f))
	}
}

func colorsFieldFilterHandler(handler func([]string) (leaf, error)) FieldFilterHandler {
	return func(colorString string) (leaf, error) {
		var colors []string

		for _, value := range colorString {
			colors = append(colors, string(value))
		}

		return handler(colors)
	}
}

func colorsEQ() FieldFilterHandler {
	return colorsFieldFilterHandler(func(colors []string) (leaf, error) {
		colorQuery := &colorQuery{
			MustInclude: colors,
			MustExclude: colorsNotInQuery(colors),
		}
		return &basicLeaf{predicator: colorQuery.toPredicator()}, nil
	})
}

func colorsLT() FieldFilterHandler {
	return colorsFieldFilterHandler(func(colors []string) (leaf, error) {
		colorQuery := &colorQuery{
			MustExclude: colorsNotInQuery(colors),
		}
		return &basicLeaf{predicator: colorQuery.toPredicator()}, nil
	})
}

func colorsGT() FieldFilterHandler {
	return colorsFieldFilterHandler(func(colors []string) (leaf, error) {
		colorQuery := &colorQuery{
			MustInclude: colors,
		}
		return &basicLeaf{predicator: colorQuery.toPredicator()}, nil
	})
}

var DefaultParser = &Parser{
	Fields: []FieldFilter{
		{
			Name: "name",
			Handlers: map[operator]FieldFilterHandler{
				opEQ: func(value string) (leaf, error) {
					return &basicLeaf{predicator: card.Name(value)}, nil
				},
			},
		},
		{
			Name:    "oracle",
			Aliases: []string{"o", "text"},
			Handlers: map[operator]FieldFilterHandler{
				opEQ: func(value string) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.OracleTextContainsFold(value))}, nil
				},
			},
		},
		{
			Name: "cmc",
			Handlers: map[operator]FieldFilterHandler{
				opEQ: floatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcEQ(value))}, nil
				}),
				opLT: floatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcLT(value))}, nil
				}),
				opLE: floatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcLTE(value))}, nil
				}),
				opGT: floatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcGT(value))}, nil
				}),
				opGE: floatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcGTE(value))}, nil
				}),
			},
		},
		{
			Name:    "colors",
			Aliases: []string{"c"},
			Handlers: map[operator]FieldFilterHandler{
				opEQ: colorsEQ(),
				opLT: colorsLT(),
				opLE: colorsLT(),
				opGT: colorsGT(),
				opGE: colorsGT(),
			},
		},
	},
}

// ParseQuery parses a query string and returns a node that can be converted to a bones predicate.
// This is the main entry point for the ql package.
func ParseQuery(query string) (node, error) {
	tokens, err := LexString(query)
	if err != nil {
		return nil, fmt.Errorf("failed to lex query: %w", err)
	}

	return DefaultParser.ParseTokens(tokens)
}
