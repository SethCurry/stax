package stax

// ColorField stores a set of colors for a card, like its colors
// or its color identity.
type ColorField uint8

func (c *ColorField) setColor(color *Color, on bool) {
	if on {
		*c |= color.ColorField()
	} else {
		*c &^= color.ColorField()
	}
}

func (c *ColorField) SetWhite(w bool) {
	c.setColor(ColorWhite, w)
}

func (c ColorField) hasColor(color *Color) bool {
	return c&color.ColorField() == color.ColorField()
}

func (c ColorField) HasWhite() bool {
	return c.hasColor(ColorWhite)
}

func (c *ColorField) SetBlue(w bool) {
	c.setColor(ColorBlue, w)
}

func (c ColorField) HasBlue() bool {
	return c.hasColor(ColorBlue)
}

func (c *ColorField) SetBlack(w bool) {
	c.setColor(ColorBlack, w)
}

func (c ColorField) HasBlack() bool {
	return c.hasColor(ColorBlack)
}

func (c *ColorField) SetRed(w bool) {
	c.setColor(ColorRed, w)
}

func (c ColorField) HasRed() bool {
	return c.hasColor(ColorRed)
}

func (c *ColorField) SetGreen(w bool) {
	c.setColor(ColorGreen, w)
}

func (c ColorField) HasGreen() bool {
	return c.hasColor(ColorGreen)
}

type Color struct {
	name       string
	char       string
	colorField ColorField
}

func (c Color) Name() string {
	return c.name
}

func (c Color) Char() string {
	return c.char
}

func (c Color) ColorField() ColorField {
	return c.colorField
}

var (
	ColorRed *Color = &Color{
		name:       "Red",
		char:       "R",
		colorField: ColorField(0b00010),
	}
	ColorBlue *Color = &Color{
		name:       "Blue",
		char:       "U",
		colorField: ColorField(0b01000),
	}
	ColorBlack *Color = &Color{
		name:       "Black",
		char:       "B",
		colorField: ColorField(0b00100),
	}
	ColorWhite *Color = &Color{
		name:       "White",
		char:       "W",
		colorField: ColorField(0b10000),
	}
	ColorGreen *Color = &Color{
		name:       "Green",
		char:       "G",
		colorField: ColorField(0b00001),
	}
)
