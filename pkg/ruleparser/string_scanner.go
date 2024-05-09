package ruleparser

func NewStringScanner(data string) *StringScanner {
	return &StringScanner{
		index: 0,
		data:  data,
	}
}

type StringScanner struct {
	index int
	data  string
}

func (s *StringScanner) Next() (rune, bool) {
	if s.index >= len(s.data) {
		return 0, false
	}

	c := rune(s.data[s.index])
	s.index++

	return c, true
}

func (s *StringScanner) ReadUntil(c []rune) string {
	acc := ""

	for s.index < len(s.data) {
		for _, r := range c {
			if rune(s.data[s.index]) == r {
				return acc
			}
		}

		acc += string(s.data[s.index])
		s.index++
	}

	return acc
}
