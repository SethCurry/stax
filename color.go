package stax

import (
	"fmt"
	"strings"
)

func NewUnknownColorError(symbol string) *UnknownColorError {
	return &UnknownColorError{
		Symbol: symbol,
	}
}

type UnknownColorError struct {
	Symbol string
}

func (e *UnknownColorError) Error() string {
	return fmt.Sprintf("unknown color: %s", e.Symbol)
}

type Color struct {
	Letter string `json:"letter"`
	Name   string `json:"name"`
}

func (c Color) Symbol() string {
	return fmt.Sprintf("{%s}", c.Letter)
}

func (c *Color) UnmarshalText(text []byte) error {
	textStr := strings.Trim(string(text), "{}")
	if len(textStr) == 0 {
		return NewUnknownColorError(textStr)
	}

	for _, color := range AllColors() {
		if color.Letter == textStr[0:1] {
			*c = color

			return nil
		}
	}

	return NewUnknownColorError(textStr[0:1])
}

func ColorRed() Color {
	return Color{
		Letter: "R",
		Name:   "Red",
	}
}

func ColorBlue() Color {
	return Color{
		Letter: "U",
		Name:   "Blue",
	}
}

func ColorBlack() Color {
	return Color{
		Letter: "B",
		Name:   "Black",
	}
}

func ColorGreen() Color {
	return Color{
		Letter: "G",
		Name:   "Green",
	}
}

func ColorWhite() Color {
	return Color{
		Letter: "W",
		Name:   "White",
	}
}

func AllColors() []Color {
	return []Color{ColorRed(), ColorBlue(), ColorBlack(), ColorGreen(), ColorWhite()}
}
