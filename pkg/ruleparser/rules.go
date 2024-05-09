package ruleparser

import (
	"bytes"
)

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

func scanLines(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
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
