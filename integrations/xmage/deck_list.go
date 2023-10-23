package xmage

import (
	"fmt"
	"strings"
)

type DeckCard struct {
	Name            string
	SetCode         string
	CollectorNumber string
	Quantity        int
}

func (d DeckCard) SetReference() string {
	return "[" + d.SetCode + ":" + d.CollectorNumber + "]"
}

func (d DeckCard) MarshalDck() string {
	return fmt.Sprintf("%d %s %s", d.Quantity, d.SetReference(), d.Name)
}

// NewDeckList creates a new DeckList.
func NewDeckList(name string, cards []DeckCard, sideboard []DeckCard) DeckList {
	return DeckList{
		Name:      name,
		Cards:     cards,
		Sideboard: sideboard,
	}
}

// DeckList represents a deck list in xmage.
type DeckList struct {
	Name      string
	Cards     []DeckCard
	Sideboard []DeckCard
}

func (d DeckList) MarshalDck() string {
	builder := strings.Builder{}

	if d.Name != "" {
		builder.WriteString(fmt.Sprintf("NAME:%s\n", d.Name))
	}

	for _, card := range d.Cards {
		builder.WriteString(card.MarshalDck() + "\n")
	}

	if len(d.Sideboard) > 0 {
		builder.WriteString("SB: ")
		for idx, card := range d.Sideboard {
			marsh := card.MarshalDck()

			if idx != 0 {
				marsh = " " + marsh
			}

			builder.WriteString(marsh)
		}
		builder.WriteString("\n")
	}

	builder.WriteString("LAYOUT MAIN:(1,1)(CMC,false,50)|(")
	for idx, v := range d.Cards {
		ref := v.SetReference()
		if idx != 0 {
			ref = "," + ref
		}
		builder.WriteString(ref)
	}
	builder.WriteString(")\n")

	if len(d.Sideboard) > 0 {
		builder.WriteString("LAYOUT SIDEBOARD:(1,1)(CMC,false,50)|(")
		for idx, v := range d.Sideboard {
			ref := v.SetReference()
			if idx != 0 {
				ref = "," + ref
			}
			builder.WriteString(ref)
		}
		builder.WriteString(")\n")
	}

	builder.WriteString("\n")
	return builder.String()
}
