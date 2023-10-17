package stax

type DeckListCard struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type DeckList struct {
	Name       string         `json:"name"`
	Primer     string         `json:"primer"`
	Commander  string         `json:"commander"`
	PowerLevel int            `json:"power_level"`
	Cards      []DeckListCard `json:"cards"`
}
