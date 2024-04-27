package scryfall

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrUnknownSource = errors.New("unknown source")

// Source is an enum representing the source of some information.
// See AllSources() for all possible values.
type Source string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *Source) UnmarshalText(txt []byte) error {
	allSources := AllSources()
	asSource := Source(string(txt))

	for _, source := range allSources {
		if source == asSource {
			*s = source

			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrUnknownSource, string(txt))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Source) UnmarshalJSON(txt []byte) error {
	var unmarshed string

	err := json.Unmarshal(txt, &unmarshed)
	if err != nil {
		return fmt.Errorf("failed to unmarshal source: %w", err)
	}

	return s.UnmarshalText([]byte(unmarshed))
}

// MarshalText implements the encoding.TextMarshaler interface.
func (s Source) MarshalText() ([]byte, error) {
	return []byte(string(s)), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (s Source) MarshalJSON() ([]byte, error) {
	return s.MarshalText()
}

const (
	// SourceWOTC indicates rulings that come directly
	// from Wizards of the Coast.
	SourceWOTC = Source("wotc")

	// SourceScryfall indicates rulings that were added by Scryfall.
	SourceScryfall = Source("scryfall")
)

// AllSources returns all possible values of Source.
func AllSources() []Source {
	return []Source{SourceWOTC, SourceScryfall}
}
