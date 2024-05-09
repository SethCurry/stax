package ruleparser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func NewSubrule(contents []*ContentElement) *Subrule {
	return &Subrule{
		Contents: contents,
		Examples: [][]*ContentElement{},
	}
}

type Subrule struct {
	Contents []*ContentElement   `json:"contents"`
	Examples [][]*ContentElement `json:"examples"`
}

func (s *Subrule) AddExample(example []*ContentElement) {
	s.Examples = append(s.Examples, example)
}

func (s *Subrule) AddToContents(text string) {
	s.Contents = append(s.Contents, parseContent(text)...)
}

func isSubrule(line string) bool {
	re := regexp.MustCompile(`^[0-9]{3}\.[0-9]+[a-z]\.? .*`)

	return re.MatchString(line)
}

var ErrInvalidSubruleLine = errors.New("invalid subrule line")

func parseSubruleLine(line string) (string, *Subrule, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("%w: not enough spaces to be valid", ErrInvalidSubruleLine)
	}

	ruleNumberWithDot := strings.TrimSuffix(parts[0], ".")
	ruleNumberParts := strings.Split(ruleNumberWithDot, ".")

	if len(ruleNumberParts) != 2 {
		return "", nil, ErrInvalidRuleNumber
	}

	ruleLetter := ruleNumberParts[1][len(ruleNumberParts[1])-1]
	text := strings.Join(parts[1:], " ")

	return string(ruleLetter), NewSubrule(parseContent(text)), nil
}
