package ruleparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

var ErrInvalidSubsectionLine = errors.New("invalid subsection line")

func parseSubsectionLine(line string) (int, *Subsection, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return 0, nil, fmt.Errorf("%w: not enough spaces to be a valid subsection line", ErrInvalidSubsectionLine)
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

func isSubsection(line string) bool {
	re := regexp.MustCompile(`^[0-9]{3}\. .*`)

	return re.MatchString(line)
}
