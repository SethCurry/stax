// Package console contains utilities for communicating with users via the command line.
package console

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/SethCurry/stax/pkg/scryfall"
)

// NewScryfallCardTable returns a new ScryfallCardTable that writes to the given io.Writer.
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

// writeEmptyRow writes an empty row to the table, (i.e. just 3 tabs separated by pipes, with a newline).
// This is useful when leaving empty rows between multi-line records in a table.
func (t *ScryfallCardTable) writeEmptyRow() error {
	_, err := t.writer.Write([]byte("\t|\t|\t|\n"))

	return err
}

func (t *ScryfallCardTable) Write(card *scryfall.Card) error {
	// double-sided or split cards have their name set to "name1 // name2"
	nameSplit := strings.Split(card.Name, "//")

	// ditto for types
	typeSplit := strings.Split(card.TypeLine, "//")

	name := nameSplit[0]
	cardType := typeSplit[0]

	newLineAfter := false

	if len(nameSplit) > 1 {
		newLineAfter = true

		err := t.writeEmptyRow()
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
		err = t.writeEmptyRow()
		if err != nil {
			return err
		}
	}

	return nil
}

// Flush tells the ScryfallCardTable to write any buffered data
// to disk.
func (t *ScryfallCardTable) Flush() {
	t.writer.Flush()
}

func ScryfallCardJSON(card *scryfall.Card) error {
	marshalled, err := json.MarshalIndent(card, "  ", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal card: %w", err)
	}

	os.Stdout.Write(marshalled)
	os.Stdout.Write([]byte(",\n"))

	defer os.Stdout.Write([]byte("]\n"))

	return nil
}
