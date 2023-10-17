package mtgrules

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ContentType string

func ContentText() ContentType {
	return ContentType("text")
}

func ContentReference() ContentType {
	return ContentType("reference")
}

func ContentSymbol() ContentType {
	return ContentType("symbol")
}

func NewRules() *Rules {
	return &Rules{
		Sections: make(map[int]*Section),
		Glossary: make(map[string]string),
	}
}

type Rules struct {
	Sections map[int]*Section  `json:"sections"`
	Glossary map[string]string `json:"glossary"`
}

func NewSection(name string) *Section {
	return &Section{
		Name:        name,
		Subsections: make(map[int]*Subsection),
	}
}

type Section struct {
	Name        string              `json:"name"`
	Subsections map[int]*Subsection `json:"subsections"`
}

func NewSubsection(name string) *Subsection {
	return &Subsection{
		Name:  name,
		Rules: make(map[int]*Rule),
	}
}

type Subsection struct {
	Name  string        `json:"name"`
	Rules map[int]*Rule `json:"rules"`
}

func NewRule(contents []*ContentElement) *Rule {
	return &Rule{
		Contents: contents,
		Subrules: make(map[string]*Subrule),
		Examples: [][]*ContentElement{},
	}
}

type Rule struct {
	Subrules map[string]*Subrule `json:"subrules"`
	Contents []*ContentElement   `json:"contents"`
	Examples [][]*ContentElement `json:"examples"`
}

func (r *Rule) AddExample(example []*ContentElement) {
	r.Examples = append(r.Examples, example)
}

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

type Exampler interface {
	AddExample(example []*ContentElement)
}

type ContentElement struct {
	Type  ContentType `json:"type"`
	Value string      `json:"value"`
}

func scanLines(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func isSection(line string) (bool, error) {
	re, err := regexp.Compile(`^[0-9]\. .*`)
	if err != nil {
		return false, err
	}

	return re.MatchString(line), nil
}

func isSubsection(line string) (bool, error) {
	re, err := regexp.Compile(`^[0-9]{3}\. .*`)
	if err != nil {
		return false, err
	}

	return re.MatchString(line), nil
}

func isRule(line string) (bool, error) {
	re, err := regexp.Compile(`^[0-9]{3}\.[0-9]+\. .*`)
	if err != nil {
		return false, err
	}
	return re.MatchString(line), nil
}

func isSubrule(line string) (bool, error) {
	re, err := regexp.Compile(`^[0-9]{3}\.[0-9]+[a-z] .*`)
	if err != nil {
		return false, err
	}
	return re.MatchString(line), nil
}

func isExample(line string) (bool, error) {
	re, err := regexp.Compile(`^Example: .*`)
	if err != nil {
		return false, err
	}

	return re.MatchString(line), nil
}

func ParseRulesFile(path string) (*Rules, error) {
	fileDescriptor, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fileDescriptor.Close()

	return ParseRules(fileDescriptor)
}

func parseSectionLine(line string) (int, *Section, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return 0, nil, fmt.Errorf("invalid section line: %s", line)
	}

	sectionNumber := parts[0]
	sectionNumber = strings.TrimSuffix(sectionNumber, ".")
	asInt, err := strconv.Atoi(sectionNumber)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse section number \"%s\" as integer: %w", sectionNumber, err)
	}
	sectionName := strings.Join(parts[1:], " ")

	return asInt, NewSection(sectionName), nil
}

func parseSubsectionLine(line string) (int, *Subsection, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return 0, nil, fmt.Errorf("invalid subsection line: %s", line)
	}

	sectionNumber := parts[0]
	sectionNumber = strings.TrimSuffix(sectionNumber, ".")
	asInt, err := strconv.Atoi(sectionNumber)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse subsection number \"%s\" as integer: %w", sectionNumber, err)
	}
	sectionName := strings.Join(parts[1:], " ")

	return asInt, NewSubsection(sectionName), nil
}

func parseRuleLine(line string) (int, *Rule, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return 0, nil, fmt.Errorf("invalid rule line: %s", line)
	}

	ruleNumberWithDot := strings.TrimSuffix(parts[0], ".")
	ruleNumberParts := strings.Split(ruleNumberWithDot, ".")
	if len(ruleNumberParts) != 2 {
		return 0, nil, fmt.Errorf("invalid rule number: %s", ruleNumberWithDot)
	}

	ruleNumberString := ruleNumberParts[1]
	ruleNumber, err := strconv.Atoi(ruleNumberString)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid rule number integer: %w", err)
	}

	text := strings.Join(parts[1:], " ")

	return ruleNumber, NewRule(parseContent(text)), nil
}

func parseExample(line string) ([]*ContentElement, error) {
	line = line[len("Example: "):]
	return parseContent(line), nil
}

func parseSubruleLine(line string) (string, *Subrule, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("invalid rule line: %s", line)
	}

	ruleNumberWithDot := strings.TrimSuffix(parts[0], ".")
	ruleNumberParts := strings.Split(ruleNumberWithDot, ".")
	if len(ruleNumberParts) != 2 {
		return "", nil, fmt.Errorf("invalid rule number: %s", ruleNumberWithDot)
	}

	ruleLetter := ruleNumberParts[1][len(ruleNumberParts[1])-1]
	text := strings.Join(parts[1:], " ")

	return string(ruleLetter), NewSubrule(parseContent(text)), nil
}

func NewStringScanner(data string) *StringScanner {
	return &StringScanner{
		data: data,
	}
}

type StringScanner struct {
	index int
	data  string
}

func (s *StringScanner) Next() (rune, bool) {
	if s.index >= len(s.data) {
		return 0, false
	}

	c := rune(s.data[s.index])
	s.index++
	return c, true
}

func (s *StringScanner) ReadUntil(c []rune) string {
	acc := ""
	for s.index < len(s.data) {
		for _, r := range c {
			if rune(s.data[s.index]) == r {
				return acc
			}
		}
		acc += string(s.data[s.index])
		s.index++
	}

	return acc
}

func parseContent(content string) []*ContentElement {
	refRegex := regexp.MustCompile(`[0-9]{3}\.[0-9]+[a-z]?`)
	elements := []*ContentElement{}
	acc := ""

	scanner := NewStringScanner(content)

	for {
		c, ok := scanner.Next()
		if !ok {
			break
		}

		switch c {
		case '{':
			if len(acc) > 0 {
				elements = append(elements, &ContentElement{
					Type:  ContentText(),
					Value: acc,
				})
				acc = ""
			}
			value := scanner.ReadUntil([]rune{'}'})
			elements = append(elements, &ContentElement{
				Type:  ContentSymbol(),
				Value: value,
			})
			scanner.Next()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			maybeRef := string(c) + scanner.ReadUntil([]rune{' ', ')'})
			maybeRef = strings.TrimSuffix(maybeRef, ".")
			if refRegex.MatchString(maybeRef) {
				if len(acc) > 0 {
					elements = append(elements, &ContentElement{
						Type:  ContentText(),
						Value: acc,
					})
					acc = ""
				}
				elements = append(elements, &ContentElement{
					Type:  ContentReference(),
					Value: maybeRef,
				})
			} else {
				acc += string(maybeRef)
			}
		default:
			acc += string(c)
		}
	}

	if len(acc) > 0 {
		elements = append(elements, &ContentElement{
			Type:  ContentText(),
			Value: acc,
		})
	}

	return elements
}

func ParseRules(reader io.Reader) (*Rules, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLines)

	inRules := false
	inGlossary := false

	rules := NewRules()

	var currentSection *Section
	var currentSubsection *Subsection
	var currentRule *Rule
	var lastExampler Exampler

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "Glossary" {
			if inRules {
				inGlossary = true
			}
			inRules = true
			continue
		}

		if !inRules {
			continue
		}

		if inGlossary {
			continue
		}

		isSect, err := isSection(line)
		if err != nil {
			return nil, err
		}
		if isSect {
			num, sect, err := parseSectionLine(line)
			if err != nil {
				return nil, err
			}
			currentSection = sect
			rules.Sections[num] = sect

			continue
		}

		isSubsect, err := isSubsection(line)
		if err != nil {
			return nil, err
		}
		if isSubsect {
			num, subsect, err := parseSubsectionLine(line)
			if err != nil {
				return nil, err
			}
			currentSubsection = subsect
			currentSection.Subsections[num] = subsect
			continue
		}

		isRuleLine, err := isRule(line)
		if err != nil {
			return nil, err
		}
		if isRuleLine {
			num, r, err := parseRuleLine(line)
			if err != nil {
				return nil, err
			}
			currentRule = r
			lastExampler = r
			currentSubsection.Rules[num] = r
			continue
		}

		isSubruleLine, err := isSubrule(line)
		if err != nil {
			return nil, err
		}
		if isSubruleLine {
			l, r, err := parseSubruleLine(line)
			if err != nil {
				return nil, err
			}
			currentRule.Subrules[l] = r
			lastExampler = r
			continue
		}

		isExampleLine, err := isExample(line)
		if err != nil {
			return nil, err
		}
		if isExampleLine {
			example, err := parseExample(line)
			if err != nil {
				return nil, err
			}
			lastExampler.AddExample(example)
			continue
		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}
	return rules, nil
}
