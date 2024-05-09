package ruleparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

func isSection(line string) bool {
	re := regexp.MustCompile(`^[0-9]\. .*`)

	return re.MatchString(line)
}

var ErrInvalidSectionLine = errors.New("invalid section line")

func parseSectionLine(line string) (int, *Section, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return 0, nil, fmt.Errorf("%w: not enough spaces to be a valid section line", ErrInvalidSectionLine)
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
