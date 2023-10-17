package stax

type DeckPercentageByType struct{}

type DeckPercentageByColor struct{}

type DeckReport struct {
	AverageCMC          float32
	AverageCMCWithLands float32
	CardCount           int
	PercentageByType    DeckPercentageByType
	PercentageByColor   DeckPercentageByColor
}
