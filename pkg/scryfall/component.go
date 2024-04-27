package scryfall

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrUnknownComponent = errors.New("unknown component")

// Component represents the "component" field that
// is present on some Scryfall responses.
// They are used to express some relationships between cards,
// like combo pieces or meld parts/results.
type Component string

const (
	// ComponentToken indicates that the card is a token.
	// This can be useful because tokens do not have unique names.
	ComponentToken Component = "token"

	// ComponentComboPiece indicates that the card is a combo piece.
	ComponentComboPiece Component = "combo_piece"

	// ComponentMeldPart indicates that the card is a meld part.
	ComponentMeldPart Component = "meld_part"

	// ComponentMeldResult indicates that the card is a meld result.
	ComponentMeldResult Component = "meld_result"
)

// String converts the component back into a string.
// This is useful for displaying it to the user, and for implementing
// fmt.Stringer.
func (c Component) String() string {
	return string(c)
}

func (c Component) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c Component) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(c.String())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal component: %w", err)
	}

	return marshalled, nil
}

func (c *Component) UnmarshalText(txt []byte) error {
	allComponents := AllComponents()
	asStr := string(txt)

	for _, component := range allComponents {
		if component.String() == asStr {
			*c = component

			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrUnknownComponent, string(txt))
}

func (c *Component) UnmarshalJSON(txt []byte) error {
	var unmarshed string

	err := json.Unmarshal(txt, &unmarshed)
	if err != nil {
		return fmt.Errorf("failed to unmarshal component: %w", err)
	}

	return c.UnmarshalText([]byte(unmarshed))
}

// AllComponents returns a slice of all of the known components.
func AllComponents() []Component {
	return []Component{
		ComponentToken,
		ComponentComboPiece,
		ComponentMeldPart,
		ComponentMeldResult,
	}
}
