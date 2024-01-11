package scryfall

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrUnknownLegality = errors.New("unknown legality")

// Legality is an enum representing the legality of a card in a format.
// See AllLegalities() for all possible values.
type Legality string

// String returns the legality as a string.
func (l Legality) String() string {
	return string(l)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (l *Legality) UnmarshalText(txt []byte) error {
	allLegalities := AllLegalities()

	for _, legality := range allLegalities {
		if legality == Legality(string(txt)) {
			*l = legality

			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrUnknownLegality, string(txt))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *Legality) UnmarshalJSON(txt []byte) error {
	var unmarshed string

	err := json.Unmarshal(txt, &unmarshed)
	if err != nil {
		return fmt.Errorf("failed to unmarshal legality: %w", err)
	}

	return l.UnmarshalText([]byte(unmarshed))
}

// MarshalText implements the encoding.TextMarshaler interface.
func (l Legality) MarshalText() ([]byte, error) {
	return []byte(string(l)), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (l Legality) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(l.String())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal legality: %w", err)
	}

	return marshalled, nil
}

// LegalityNotLegal returns a Legality representing not legal, e.g. cards out of rotation.
func LegalityNotLegal() Legality {
	return Legality("not_legal")
}

// LegalityLegal returns a Legality representing a card that is legal in a format.
func LegalityLegal() Legality {
	return Legality("legal")
}

// LegalityRestricted returns a Legality representing a card that is restricted in a format.
func LegalityRestricted() Legality {
	return Legality("restricted")
}

// LegalityBanned returns a Legality representing a card that is specifically banned.
func LegalityBanned() Legality {
	return Legality("banned")
}

func AllLegalities() []Legality {
	return []Legality{
		LegalityNotLegal(),
		LegalityLegal(),
		LegalityRestricted(),
		LegalityBanned(),
	}
}
