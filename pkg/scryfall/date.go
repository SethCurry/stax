package scryfall

import (
	"encoding/json"
	"fmt"
	"time"
)

// Date is a type alias for time.Time that supports marshalling/unmarshalling
// text and JSON in the format YYYY-MM-DD.
type Date time.Time

// UnmarshalJSON implements the json.Unmarshaler interface.
// It expects the date to be in the format YYYY-MM-DD.
func (d *Date) UnmarshalJSON(txt []byte) error {
	var unmarshed string

	err := json.Unmarshal(txt, &unmarshed)
	if err != nil {
		return fmt.Errorf("failed to unmarshal date: %w", err)
	}

	return d.UnmarshalText([]byte(unmarshed))
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It expects the date to be in the format YYYY-MM-DD.
func (d *Date) UnmarshalText(txt []byte) error {
	parsedTime, err := time.Parse("2006-01-02", string(txt))
	if err != nil {
		return fmt.Errorf("date is not valid: %w", err)
	}

	*d = Date(parsedTime)

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The generated date will be in the format YYYY-MM-DD.
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// MarshalJSON implements the json.Marshaler interface.
// The generated date will be in the format YYYY-MM-DD.
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// String returns the given date as a string in the format
// YYYY-MM-DD.
func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}
