package ql

import (
	"errors"
	"strings"

	"github.com/SethCurry/scurry-go/fp"
)

type keyword string

const (
	keywordOr  = "OR"
	keywordAnd = "AND"
)

var keywords = []keyword{
	keywordOr,
	keywordAnd,
}

// TokenFamily represents the general category of a token,
// such as a literal, keyword, or operator.
type TokenFamily string

const (
	// FamilyOperator represents an operator token such as >, <, =, !=, etc.
	FamilyOperator TokenFamily = "operator"

	// FamilyLiteral represents a literal token such as a string or number.
	FamilyLiteral TokenFamily = "literal"

	// FamilyKeyword represents a keyword token such as OR or AND.
	FamilyKeyword TokenFamily = "keyword"

	// FamilyParen represents a parenthesis token such as ( or ).
	FamilyParen TokenFamily = "paren"
)

// Token represents a single token in the query.
type Token struct {
	// Family represents the general category of the token.
	Family TokenFamily

	// Value represents the actual value of the token from the query.
	Value string
}

// isKeyword checks if a given string is a keyword, case-insensitively.
// It returns true if the string is a keyword, false otherwise.
func isKeyword(item string) bool {
	for _, v := range keywords {
		if strings.ToUpper(item) == string(v) {
			return true
		}
	}

	return false
}

// tokenLiteralsToKeywords accepts a slice of tokens and returns a new slice
// with all literal tokens whose value is a keyword converted to a keyword token.
func tokenLiteralsToKeywords(tokens []Token) []Token {
	return fp.Map(func(t Token) Token {
		if isKeyword(t.Value) {
			t.Value = strings.ToUpper(t.Value)
			t.Family = FamilyKeyword
		}

		return t
	}, tokens)
}

func lex(t *lexReader) ([]Token, error) {
	// TODO attach line/column info to tokens so it can be propagated in errors
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
			case ">", "<":
				// if the next token is an =, then it's a combined operator
				if next, ok := t.peek(); ok && next == '=' {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  nextItem + "=",
					})
					_, done = t.next()
				} else {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  nextItem,
					})
				}
			case "=", ":":
				// Same thing as < and >, if the next token is a > or <, then it's a combined operator
				// This makes users able to use them in either order, i.e. >= or =>
				if next, ok := t.peek(); ok && (next == '>' || next == '<') {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  string(next) + "=",
					})

					// consume the next character, since it's part of the operator
					_, done = t.next()
				} else {
					ret = append(ret, Token{
						Family: FamilyOperator,
						Value:  "=",
					})
				}
			case "(", ")":
				ret = append(ret, Token{
					Family: FamilyParen,
					Value:  nextItem,
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

// LexString takes a string and returns a slice of tokens.
// You will typically want to just use ParseQuery, but this can be useful
// if you want to inspect or manipulate the tokens before parsing.
func LexString(input string) ([]Token, error) {
	proc := &lexReader{
		query: []rune(input),
	}

	return lex(proc)
}
