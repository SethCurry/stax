package ql

import (
	"errors"
	"fmt"
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

// leaf is a terminal node in the parse tree.
// These are typically basicLeaf nodes implementing a filter like
// "name:Static"
type leaf interface {
	Predicate() predicate.Card
}

// node is a single node in the parse tree.
// These are typically AND or OR nodes.
type node interface {
	leaf
	Left() leaf
	Right() leaf
	SetLeft(leaf)
	SetRight(leaf)
}

// logicNode is a node that holds a logic operator (AND or OR) and two children.
type logicNode struct {
	predicator func(...predicate.Card) predicate.Card
	left       leaf
	right      leaf
}

// Predicate returns the predicate for the node,
// such as ANDing or ORing two predicates together.
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

// Left returns the left child of the node.
// This is required to satisfy the node interface.
func (l *logicNode) Left() leaf {
	return l.left
}

// Right returns the right child of the node.
// This is required to satisfy the node interface.
func (l *logicNode) Right() leaf {
	return l.right
}

// SetLeft sets the left child of the node.
// This is required to satisfy the node interface.
func (l *logicNode) SetLeft(left leaf) {
	l.left = left
}

// SetRight sets the right child of the node.
// This is required to satisfy the node interface.
func (l *logicNode) SetRight(right leaf) {
	l.right = right
}

// newAndNode creates a new AND node.
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

// newOrNode creates a new OR node.
func newOrNode(left, right leaf) *logicNode {
	return &logicNode{
		predicator: func(cards ...predicate.Card) predicate.Card {
			return card.Or(cards...)
		},
		left:  left,
		right: right,
	}
}

// basicLeaf is a leaf node that holds a single predicate.
type basicLeaf struct {
	predicator predicate.Card
}

// Predicate returns the predicate for the leaf node.
// This is required to satisfy the leaf interface.
func (l *basicLeaf) Predicate() predicate.Card {
	return l.predicator
}

// newTokenReader creates a new token reader for a slice of tokens.
func newTokenReader(tokens []Token) *tokenReader {
	return &tokenReader{
		tokens: tokens,
		index:  0,
	}
}

// tokenReader is a helper struct for reading tokens from a slice.
// It allows iterating over the tokens non-linearly (i.e. by calling next() multiple times).
type tokenReader struct {
	tokens []Token
	index  int
}

// next returns the next token in the slice and increments the index to move to the next token.
// It returns nil, false if there are no more tokens to read.
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

// hasMore returns true if there are more tokens to read,
// or false if the index is out of bounds.
func (r *tokenReader) hasMore() bool {
	return r.index < len(r.tokens)
}

// Parser holds a list of fields that have been registered and will be used to parse queries.
// You can add custom fields to the parser by creating a new *Parser and adding FieldFilters to it.
type Parser struct {
	Fields []FieldFilter
}

// AddField adds a field to the parser.
// Used to add custom fields to the parser.
func (p *Parser) AddField(field FieldFilter) {
	p.Fields = append(p.Fields, field)
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

// ParseQuery parses a query string and returns a node that can be converted to a bones predicate.
func (p *Parser) ParseQuery(query string) (node, error) {
	tokens, err := LexString(query)
	if err != nil {
		return nil, fmt.Errorf("failed to lex query: %w", err)
	}

	return p.ParseTokens(tokens)
}

// DefaultParser is a parser with the standard set of fields already registered.
// This is the parser that you typically want to use.
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
				opEQ: FloatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcEQ(value))}, nil
				}),
				opLT: FloatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcLT(value))}, nil
				}),
				opLE: FloatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcLTE(value))}, nil
				}),
				opGT: FloatFieldFilterHandler(func(value float32) (leaf, error) {
					return &basicLeaf{predicator: card.HasFacesWith(cardface.CmcGT(value))}, nil
				}),
				opGE: FloatFieldFilterHandler(func(value float32) (leaf, error) {
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
