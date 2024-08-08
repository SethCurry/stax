package ql

// lexReader is a helper struct for reading tokens from a query string.
type lexReader struct {
	index int
	query []rune
}

// next reads the next character from the query string and advances the index.
func (l *lexReader) next() (rune, bool) {
	if l.index >= len(l.query) {
		return ' ', false
	}

	char := l.query[l.index]

	l.index++

	return char, true
}

// peek returns the next character from the query string without advancing the index.
func (l *lexReader) peek() (rune, bool) {
	if l.index >= len(l.query) {
		return ' ', false
	}

	return l.query[l.index], true
}

// readUntilOneOf reads characters from the query string until one of the provided
// characters is encountered. It returns the accumulated string and a boolean
// indicating if the end of the query string was reached.
func (l *lexReader) readUntilOneOf(matches []rune) (string, bool) {
	acc := ""

	for {
		nextChar, ok := l.peek()
		if !ok {
			return acc, false
		}

		for _, m := range matches {
			if nextChar == m {
				return acc, true
			}
		}
		l.next()
		acc += string(nextChar)
	}
}

// readUntilSeparator reads characters from the query string until one of the
// provided separators is encountered. It returns the accumulated string and
// a boolean indicating if the end of the query string was reached.
func (l *lexReader) readUntilSeparator() (string, bool) {
	separators := []rune{
		' ',
		'"',
		'=',
		'>',
		'<',
		'(',
		')',
		':',
	}

	return l.readUntilOneOf(separators)
}
