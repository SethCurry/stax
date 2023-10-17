package stax

import (
	"fmt"
	"strings"
)

func NewErrInvalidManaSymbol(problem string, symbol string) *ErrInvalidManaSymbol {
	return &ErrInvalidManaSymbol{
		Problem: problem,
		Symbol:  symbol,
	}
}

type ErrInvalidManaSymbol struct {
	Problem string
	Symbol  string
}

func (e *ErrInvalidManaSymbol) Error() string {
	return fmt.Sprintf("invalid mana symbol: %s: %s", e.Problem, e.Symbol)
}

var _ Mana = &ManaColored{}

type ManaColored struct {
	Color *Color `json:"color"`
}

func (m ManaColored) Symbol() string {
	return m.Color.Symbol()
}

func (m ManaColored) Colors() []*Color {
	return []*Color{m.Color}
}

var _ Mana = &ManaColorless{}

type ManaColorless struct{}

func (m ManaColorless) Symbol() string {
	return "{C}"
}

func (m ManaColorless) Colors() []*Color {
	return []*Color{}
}

var _ Mana = &ManaGeneric{}

type ManaGeneric struct {
	Amount int `json:"amount"`
}

func (m ManaGeneric) Symbol() string {
	if m.Amount == -1 {
		return "{X}"
	} else {
		return fmt.Sprintf("{%d}", m.Amount)
	}
}

func (m ManaGeneric) Colors() []*Color {
	return []*Color{}
}

var _ Mana = &ManaPhyrexian{}

type ManaPhyrexian struct {
	Color *Color `json:"color"`
}

func (m ManaPhyrexian) Symbol() string {
	return fmt.Sprintf("{%s/P}", m.Color.Letter)
}

func (m ManaPhyrexian) Colors() []*Color {
	return []*Color{m.Color}
}

var _ Mana = &ManaHybrid{}

type ManaHybrid struct {
	Color []*Color `json:"colors"`
}

func (m ManaHybrid) Symbol() string {
	symbols := make([]string, len(m.Color))
	for i, color := range m.Color {
		symbols[i] = color.Symbol()
	}

	return fmt.Sprintf("{%s}", strings.Join(symbols, "/"))
}

func (m ManaHybrid) Colors() []*Color {
	return m.Color
}

type Mana interface {
	Symbol() string
	Colors() []*Color
}
