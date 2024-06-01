package stax

import (
	"fmt"
	"io"
)

func NewMTGODecklistWriter(writer io.Writer) *MTGODecklistWriter {
	return &MTGODecklistWriter{
		writer: writer,
	}
}

type MTGODecklistWriter struct {
	writer io.Writer
}

func (m *MTGODecklistWriter) AddCard(name string, count int) {
	m.writer.Write([]byte(fmt.Sprintf("%d %s\n", count, name)))
}
