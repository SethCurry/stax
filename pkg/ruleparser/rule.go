package ruleparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// NewRule creates a new rule.  This is typically only used during
// parsing.
func NewRule(contents []*ContentElement) *Rule {
	return &Rule{
		Contents: contents,
		Subrules: make(map[string]*Subrule),
		Examples: [][]*ContentElement{},
	}
}

// Rule represents a rule in the ruleset.  It stores the contents of itself,
// as well as children like subrules and examples.
type Rule struct {
	Subrules map[string]*Subrule `json:"subrules"`
	Contents []*ContentElement   `json:"contents"`
	Examples [][]*ContentElement `json:"examples"`
}

// AddSubrule adds a subrule to the rule.
func (r *Rule) AddExample(example []*ContentElement) {
	r.Examples = append(r.Examples, example)
}

func (r *Rule) AddToContents(text string) {
	r.Contents = append(r.Contents, parseContent(text)...)
}

func isRule(line string) bool {
	re := regexp.MustCompile(`^[0-9]{3}\.[0-9]+\.? .*`)

	return re.MatchString(line) || strings.HasPrefix(line, "901.4.All")
}

// ErrInvalidRuleLine means that the line could not be parsed as a rule.
// This typically means the line is not a rule, but not necessarily.
var ErrInvalidRuleLine = errors.New("invalid rule line")

// ErrInvalidRuleNumber indicates that the rule number could not be parsed.
// This either means it does not look like a rule number (i.e. is not formatted correctly)
// or that some part of it is not an integer.
var ErrInvalidRuleNumber = errors.New("invalid rule number")

func parseRuleLine(line string) (int, *Rule, error) {
	parts := strings.Split(line, " ")

	if len(parts) < 2 {
		return 0, nil, ErrInvalidRuleLine
	}

	if parts[0] == "901.4.All" {
		return 4, NewRule(parseContent(line[6:])), nil
	}

	ruleNumberWithDot := strings.TrimSuffix(parts[0], ".")
	ruleNumberParts := strings.Split(ruleNumberWithDot, ".")

	if len(ruleNumberParts) != 2 {
		return 0, nil, ErrInvalidRuleNumber
	}

	ruleNumberString := ruleNumberParts[1]

	ruleNumber, err := strconv.Atoi(ruleNumberString)
	if err != nil {
		return 0, nil, fmt.Errorf("%w: %w", ErrInvalidRuleNumber, err)
	}

	text := strings.Join(parts[1:], " ")

	return ruleNumber, NewRule(parseContent(text)), nil
}
