package stax

import (
	"fmt"
	"strings"
)

func NewErrUnknownColor(symbol string) *ErrUnknownColor {
	return &ErrUnknownColor{
		Symbol: symbol,
	}
}

type ErrUnknownColor struct {
	Symbol string
}

func (e *ErrUnknownColor) Error() string {
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
		return NewErrUnknownColor(textStr)
	}

	for _, color := range AllColors {
		if color.Letter == textStr[0:1] {
			*c = *color
			return nil
		}
	}

	return NewErrUnknownColor(textStr[0:1])
}

var (
	ColorRed = &Color{
		Letter: "R",
		Name:   "Red",
	}
	ColorBlue = &Color{
		Letter: "U",
		Name:   "Blue",
	}
	ColorBlack = &Color{
		Letter: "B",
		Name:   "Black",
	}
	ColorGreen = &Color{
		Letter: "G",
		Name:   "Green",
	}
	ColorWhite = &Color{
		Letter: "W",
		Name:   "White",
	}

	AllColors = []*Color{
		ColorRed,
		ColorBlue,
		ColorBlack,
		ColorGreen,
		ColorWhite,
	}
)
