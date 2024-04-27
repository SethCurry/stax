package scryfall

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrUnknownObject = errors.New("unknown object")

// Object represents the "object" field that is present in all Scryfall API objects.
type Object string

// String returns the string representation of the object.
func (o Object) String() string {
	return string(o)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (o *Object) UnmarshalText(text []byte) error {
	allObjects := AllObjects()
	asStr := string(text)

	for _, v := range allObjects {
		if v.String() == asStr {
			*o = v

			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrUnknownObject, asStr)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *Object) UnmarshalJSON(text []byte) error {
	var unmarshed string

	err := json.Unmarshal(text, &unmarshed)
	if err != nil {
		return fmt.Errorf("failed to unmarshal object: %w", err)
	}

	return o.UnmarshalText([]byte(unmarshed))
}

// MarshalText implements the encoding.TextMarshaler interface.
func (o Object) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (o Object) MarshalJSON() ([]byte, error) {
	return o.MarshalText()
}

const (
	ObjectList = Object("list")

	// ObjectBulkData represents a bulk data source.
	ObjectBulkData    = Object("bulk_data")
	ObjectRuling      = Object("ruling")
	ObjectCard        = Object("card")
	ObjectRelatedCard = Object("related_card")
	ObjectCatalog     = Object("catalog")
)

// ObjectList returns a list of all valid values of Object.
func AllObjects() []Object {
	return []Object{
		ObjectCard,
		ObjectRelatedCard,
		ObjectRuling,
		ObjectBulkData,
		ObjectList,
	}
}
