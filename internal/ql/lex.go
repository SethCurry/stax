package ql

import (
	"errors"
	"strings"

	"github.com/SethCurry/scurry-go/fp"
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
		if strings.ToUpper(item) == string(v) {
			return true
		}
	}

	return false
}

func tokenLiteralsToKeywords(tokens []Token) []Token {
	return fp.Map(func(t Token) Token {
		if isKeyword(t.Value) {
			t.Value = strings.ToUpper(t.Value)
			t.Family = FamilyKeyword
		}

		return t
	}, tokens)
}

type lexReader struct {
	index int
	query []rune
}

func (l *lexReader) next() (rune, bool) {
	if l.index >= len(l.query) {
		return ' ', false
	}

	char := l.query[l.index]

	l.index++

	return char, true
}

func (l *lexReader) peek() (rune, bool) {
	if l.index >= len(l.query) {
		return ' ', false
	}

	return l.query[l.index], true
}

func (l *lexReader) readUntilOneOf(matches []rune) (string, bool) {
	acc := ""

	for {
		nextChar, ok := l.peek()
		if !ok {
			return acc, false
		}

		for _, m := range matches {
			if nextChar == m {
				return acc, true
			}
		}
		l.next()
		acc += string(nextChar)
	}
}

func (l *lexReader) readUntilSeparator() (string, bool) {
	separators := []rune{
		' ',
		'"',
		'=',
		'>',
		'<',
		'(',
		')',
		':',
	}

	return l.readUntilOneOf(separators)
}

func lex(t *lexReader) ([]Token, error) {
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
			case "=", ":":
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

func Lex(input string) ([]Token, error) {
	proc := &lexReader{
		query: []rune(input),
	}

	return lex(proc)
}
