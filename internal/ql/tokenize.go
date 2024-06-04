package ql

import (
	"errors"

	"github.com/SethCurry/stax/internal/common"
)

type Keyword string

const (
	KeywordOr  = "OR"
	KeywordAnd = "AND"
)

var Keywords = []Keyword{
	KeywordOr,
	KeywordAnd,
}

type TokenFamily string

const (
	FamilyOperator = "operator"
	FamilyLiteral  = "literal"
	FamilyKeyword  = "keyword"
	FamilyParen    = "paren"
)

type Token struct {
	Family TokenFamily
	Value  string
}

func isKeyword(item string) bool {
	for _, v := range Keywords {
		if item == string(v) {
			return true
		}
	}

	return false
}

func tokenLiteralsToKeywords(tokens []Token) []Token {
	return common.Map(tokens, func(t Token) Token {
		if isKeyword(t.Value) {
			t.Family = FamilyKeyword
		}

		return t
	})
}

type tokenizer struct {
	index int
	query []rune
}

func (t *tokenizer) next() (rune, bool) {
	if t.index >= len(t.query) {
		return ' ', false
	}

	char := t.query[t.index]

	t.index++

	return char, true
}

func (t *tokenizer) peek() (rune, bool) {
	if t.index >= len(t.query) {
		return ' ', false
	}

	return t.query[t.index], true
}

func (t *tokenizer) readUntilOneOf(matches []rune) (string, bool) {
	acc := ""

	for {
		nextChar, ok := t.peek()
		if !ok {
			return acc, false
		}

		for _, m := range matches {
			if nextChar == m {
				return acc, true
			}
		}
		t.next()
		acc += string(nextChar)
	}
}

func (t *tokenizer) readUntilSeparator() (string, bool) {
	separators := []rune{
		' ',
		'"',
		'=',
		'>',
		'<',
		'(',
		')',
	}

	return t.readUntilOneOf(separators)
}

func (t *tokenizer) run() ([]Token, error) {
	var ret []Token

	for {
		nextItem, done := t.readUntilSeparator()
		if nextItem == "" {
			nextChar, charDone := t.next()
			nextItem = string(nextChar)
			done = charDone
		}

		if nextItem != " " && nextItem != "" {
			switch nextItem {
			case ">":
				if next, ok := t.peek(); ok && next == '=' {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  ">=",
					})
					_, done = t.next()
				} else {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  ">",
					})
				}
			case "<":
				if next, ok := t.peek(); ok && next == '=' {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  "<=",
					})
					_, done = t.next()
				} else {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  "<",
					})
				}
			case "=":
				if next, ok := t.peek(); ok && (next == '>' || next == '<') {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  string(next) + "=",
					})
					_, done = t.next()
				} else {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  "=",
					})
				}
			case "(":
				ret = append(ret, Token{
					Family: FamilyParen,
					Value:  "(",
				})
			case ")":
				ret = append(ret, Token{
					Family: FamilyParen,
					Value:  ")",
				})
			case "\"":
				quoted, ok := t.readUntilOneOf([]rune{'"'})
				if !ok {
					return ret, errors.New("mismatched quotes")
				}
				_, done = t.next()

				ret = append(ret, Token{
					Family: FamilyLiteral,
					Value:  quoted,
				})
			default:
				ret = append(ret, Token{
					Family: FamilyLiteral,
					Value:  nextItem,
				})
			}
		}

		if !done {
			return tokenLiteralsToKeywords(ret), nil
		}
	}
}

func Tokenize(input string) ([]Token, error) {
	proc := tokenizer{
		query: []rune(input),
	}

	return proc.run()
}
