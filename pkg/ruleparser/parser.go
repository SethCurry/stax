package ruleparser

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func ParseFile(path string) (*Rules, error) {
	fileDescriptor, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open rules file: %w", err)
	}
	defer fileDescriptor.Close()

	return ParseRules(fileDescriptor)
}

func ParseRules(reader io.Reader) (*Rules, error) {
	pa := newParser(reader)

	return pa.execute()
}

type ParseError struct {
	Line       string
	LineNumber int
	Err        error
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("error parsing line %d: %s: %v", p.LineNumber, p.Line, p.Err)
}

func newParser(reader io.Reader) *parser {
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLines)

	return &parser{
		inRules:           false,
		inGlossary:        false,
		currentSection:    nil,
		currentSubsection: nil,
		currentRule:       nil,
		lastExampler:      nil,
		scanner:           scanner,
		rules:             NewRules(),
	}
}

type parser struct {
	// inRules tracks the state of whether the reader is current returning rules
	// as opposed to the foreword/licensing/etc that come before the actual rules.
	inRules bool

	// inGlossary tracks the state of whether the reader is currently returning the glossary.
	// Used to ignore content in the glossary.
	inGlossary bool

	// currentSection stores a pointer to the last section header that was parsed.
	// This allows subsections and rules to be added to them as that section is being parsed.
	currentSection *Section

	// currentSubsection stores a pointer to the last subsection header that was parsed.
	// This allows rules to be added to it as that subsection is being parsed.
	currentSubsection *Subsection

	// currentRule stores a pointer to the last rule header that was parsed.
	// This allows subrules to be added to it as that rule is being parsed.
	currentRule *Rule

	// lastExampler stores a pointer to the last rule or subrule that was parsed.
	// This allows examples to be added to it as that rule or subrule is being parsed.
	lastExampler Exampler

	scanner *bufio.Scanner

	rules *Rules
}

func (r *parser) parseSubsection(line string) error {
	num, subsect, err := parseSubsectionLine(line)
	if err != nil {
		return err
	}

	r.currentSubsection = subsect
	r.currentSection.Subsections[num] = subsect

	return nil
}

func (r *parser) parseSection(line string) error {
	num, sect, err := parseSectionLine(line)
	if err != nil {
		return err
	}

	r.currentSection = sect
	r.rules.Sections[num] = sect

	return nil
}

func (r *parser) parseRule(line string) error {
	num, gotRule, err := parseRuleLine(line)
	if err != nil {
		return err
	}

	r.currentRule = gotRule
	r.currentSubsection.Rules[num] = gotRule
	r.lastExampler = gotRule

	return nil
}

func (r *parser) parseSubrule(line string) error {
	l, subrule, err := parseSubruleLine(line)
	if err != nil {
		return err
	}

	r.currentRule.Subrules[l] = subrule
	r.lastExampler = subrule

	return nil
}

func (r *parser) parseExample(line string) {
	example := parseExample(line)

	r.lastExampler.AddExample(example)
}

func (r *parser) handleLine(line string, origLine string) error {
	switch {
	case isSection(line):
		return r.parseSection(line)
	case isSubsection(line):
		return r.parseSubsection(line)
	case isRule(line):
		return r.parseRule(line)
	case isSubrule(line):
		return r.parseSubrule(line)
	case isExample(line):
		r.parseExample(line)

		return nil
	case strings.HasPrefix(origLine, "    ") || strings.HasPrefix(origLine, "\n"):
		if r.lastExampler != nil {
			r.lastExampler.AddToContents(line)
		} else {
			return errors.New("unexpected content line")
		}
		return nil
	default:
		return ErrUnknownLineType
	}
}

var ErrUnknownLineType = errors.New("unknown line type")

func (r *parser) execute() (*Rules, error) {
	lineNumber := 0
	for r.scanner.Scan() {
		lineNumber++

		origLine := convertEncoding(r.scanner.Text())
		line := strings.TrimSpace(origLine)

		origLine = strings.Trim(origLine, "\n")

		if line == "" || line == "Credits" {
			continue
		}

		if line == "Glossary" {
			if r.inRules {
				r.inGlossary = true
			}

			r.inRules = true

			continue
		}

		if !r.inRules {
			continue
		}

		if r.inGlossary {
			continue
		}

		err := r.handleLine(line, origLine)
		if err != nil {
			origLineJSON, jsonErr := json.Marshal(origLine)
			if jsonErr != nil {
				return nil, &ParseError{
					Line:       origLine,
					LineNumber: lineNumber,
					Err:        err,
				}
			}
			return nil, &ParseError{
				Line:       string(origLineJSON),
				LineNumber: lineNumber,
				Err:        err,
			}
		}
	}

	if err := r.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner returned error: %w", err)
	}

	return r.rules, nil
}
