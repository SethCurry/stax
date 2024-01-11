package scryfall

import (
	"encoding/json"
	"errors"
	"fmt"
)

/*
   {
     "object": "related_card",
     "id": "c672ee51-0e77-463f-8586-191b413c44dd",
     "component": "token",
     "name": "Spirit",
     "type_line": "Token Creature — Spirit",
     "uri": "https://api.scryfall.com/cards/c672ee51-0e77-463f-8586-191b413c44dd"
   },
   {
     "object": "related_card",
     "id": "003b8c93-54d2-4f23-961e-a52d63d0a54b",
     "component": "combo_piece",
     "name": "The Restoration of Eiganjo // Architect of Restoration",
     "type_line": "Enchantment — Saga // Enchantment Creature — Fox Monk",
     "uri": "https://api.scryfall.com/cards/003b8c93-54d2-4f23-961e-a52d63d0a54b"
   }
*/

var ErrUnknownComponent = errors.New("unknown component")

type Component string

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

func ComponentToken() Component {
	return Component("token")
}

func ComponentComboPiece() Component {
	return Component("combo_piece")
}

func ComponentMeldPart() Component {
	return Component("meld_part")
}

func ComponentMeldResult() Component {
	return Component("meld_result")
}

func AllComponents() []Component {
	return []Component{
		ComponentToken(),
		ComponentComboPiece(),
		ComponentMeldPart(),
		ComponentMeldResult(),
	}
}
