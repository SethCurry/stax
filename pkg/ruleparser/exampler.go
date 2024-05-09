package ruleparser

import "regexp"

type Exampler interface {
	AddExample(example []*ContentElement)
	AddToContents(string)
}

func parseExample(line string) []*ContentElement {
	line = line[len("Example: "):]

	return parseContent(line)
}

// isExample returns whether the provided line begins an example.
func isExample(line string) bool {
	re := regexp.MustCompile(`^Example: .*`)

	return re.MatchString(line)
}
