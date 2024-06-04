package stax

// ColorField stores a set of colors for a card, like its colors
// or its color identity.
//
// Fields:
// 1 2 3 4 5
// W B B R G
// h l l e r
// i u a d e
// t e c   e
// e   k   n
type ColorField uint8

func (c *ColorField) SetColor(color *Color, on bool) {
	if on {
		*c |= color.ColorField()
	} else {
		*c &^= color.ColorField()
	}
}

func (c ColorField) HasColor(color *Color) bool {
	return c&color.ColorField() == color.ColorField()
}

func (c *ColorField) SetWhite(w bool) {
	c.SetColor(ColorWhite, w)
}

func (c ColorField) HasWhite() bool {
	return c.HasColor(ColorWhite)
}

func (c *ColorField) SetBlue(w bool) {
	c.SetColor(ColorBlue, w)
}

func (c ColorField) HasBlue() bool {
	return c.HasColor(ColorBlue)
}

func (c *ColorField) SetBlack(w bool) {
	c.SetColor(ColorBlack, w)
}

func (c ColorField) HasBlack() bool {
	return c.HasColor(ColorBlack)
}

func (c *ColorField) SetRed(w bool) {
	c.SetColor(ColorRed, w)
}

func (c ColorField) HasRed() bool {
	return c.HasColor(ColorRed)
}

func (c *ColorField) SetGreen(w bool) {
	c.SetColor(ColorGreen, w)
}

func (c ColorField) HasGreen() bool {
	return c.HasColor(ColorGreen)
}

// Color represents a single color in Magic.  There are constants declared
// for all of the colors.  See ColorRed, ColorBlue, ColorBlack, ColorGreen,
// and ColorWhite.
type Color struct {
	name       string
	char       string
	colorField ColorField
}

// Name returns the full, capitalized name of the color such as "Red" or "Blue".
func (c Color) Name() string {
	return c.name
}

// Char returns the character for the mana symbol, typically used in plain-text
// card text, such as R for {R} in mana costs.
func (c Color) Char() string {
	return c.char
}

// ColorField returns a ColorField with only this color's bit set.
func (c Color) ColorField() ColorField {
	return c.colorField
}

var (
	// ColorRed is the color Red.
	ColorRed *Color = &Color{
		name:       "Red",
		char:       "R",
		colorField: ColorField(0b00010),
	}

	// ColorBlue is the color Blue.
	ColorBlue *Color = &Color{
		name:       "Blue",
		char:       "U",
		colorField: ColorField(0b01000),
	}

	// ColorBlack is the color Black.
	ColorBlack *Color = &Color{
		name:       "Black",
		char:       "B",
		colorField: ColorField(0b00100),
	}

	// ColorWhite is the color White.
	ColorWhite *Color = &Color{
		name:       "White",
		char:       "W",
		colorField: ColorField(0b10000),
	}

	// ColorGreen is the color Green.
	ColorGreen *Color = &Color{
		name:       "Green",
		char:       "G",
		colorField: ColorField(0b00001),
	}

	// AllColors is a slice of all of the colors.
	// Useful for iterating over all colors.
	AllColors = []*Color{
		ColorRed,
		ColorBlue,
		ColorBlack,
		ColorWhite,
		ColorGreen,
	}
)
