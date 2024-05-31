package stax

// WUBRG
type ColorField uint8

func (c *ColorField) SetWhite(w bool) {
	if w {
		*c |= 0b10000
	} else {
		*c &^= 0b10000
	}
}

func (c ColorField) HasWhite() bool {
	return c&0b10000 == 0b10000
}

func (c *ColorField) SetBlue(w bool) {
	if w {
		*c |= 0b01000
	} else {
		*c &^= 0b01000
	}
}

func (c ColorField) HasBlue() bool {
	return c&0b01000 == 0b01000
}

func (c *ColorField) SetBlack(w bool) {
	if w {
		*c |= 0b00100
	} else {
		*c &^= 0b00100
	}
}

func (c ColorField) HasBlack() bool {
	return c&0b00100 == 0b00100
}

func (c *ColorField) SetRed(w bool) {
	if w {
		*c |= 0b00010
	} else {
		*c &^= 0b00010
	}
}

func (c ColorField) HasRed() bool {
	return c&0b00010 == 0b00010
}

func (c *ColorField) SetGreen(w bool) {
	if w {
		*c |= 0b00001
	} else {
		*c &^= 0b00001
	}
}

func (c ColorField) HasGreen() bool {
	return c&0b00001 == 0b00001
}
