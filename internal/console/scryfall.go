package console

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/SethCurry/stax/pkg/scryfall"
)

func NewScryfallCardTable(output io.Writer) *ScryfallCardTable {
	table := &ScryfallCardTable{
		writer: tabwriter.NewWriter(output, 0, 1, 1, ' ', 0),
	}
	table.writeHeader()

	return table
}

type ScryfallCardTable struct {
	writer *tabwriter.Writer
}

func (t *ScryfallCardTable) writeHeader() {
	fmt.Fprintln(t.writer, "Name\t| Mana Cost\t| Type\t| Set")
}

func (t *ScryfallCardTable) Write(card *scryfall.Card) error {
	nameSplit := strings.Split(card.Name, "//")
	typeSplit := strings.Split(card.TypeLine, "//")

	name := nameSplit[0]
	cardType := typeSplit[0]
	newLineAfter := false

	if len(nameSplit) > 1 {
		newLineAfter = true

		_, err := t.writer.Write([]byte("\t|\t|\t|\n"))
		if err != nil {
			return err
		}

		useType := ""

		if len(typeSplit) > 1 {
			useType = typeSplit[0] + " //"
			cardType = strings.TrimSpace(typeSplit[1])
		}

		_, err = t.writer.Write([]byte(nameSplit[0] + "//\t|\t| " + useType + "\t|\n"))
		if err != nil {
			return err
		}

		name = strings.TrimSpace(nameSplit[1])
	}

	_, err := fmt.Fprintf(t.writer, "%s\t| %s\t| %s\t| %s\n", name, card.ManaCost, cardType, card.SetName)
	if err != nil {
		return err
	}

	if newLineAfter {
		_, err = t.writer.Write([]byte("\t|\t|\t|\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *ScryfallCardTable) Flush() {
	t.writer.Flush()
}
