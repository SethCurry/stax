package ql

type lexReader struct {
	index int
	query []rune
}

func (l *lexReader) next() (rune, bool) {
	if l.index >= len(l.query) {
		return ' ', false
	}

	char := l.query[l.index]

	l.index++

	return char, true
}

func (l *lexReader) peek() (rune, bool) {
	if l.index >= len(l.query) {
		return ' ', false
	}

	return l.query[l.index], true
}

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
