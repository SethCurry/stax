package scryfall

// Ruling is a single ruling.
type Ruling struct {
	// The type of object; should be "ruling" for rulings.
	Object Object `json:"object"`

	// The card's oracle ID, set by Wizards.
	OracleID string `json:"oracle_id"`

	// The source of the ruling; should be either "wotc" or "scryfall".
	Source Source `json:"source"`

	// The date when the ruling was published.
	PublishedAt Date `json:"published_at"`

	// The actual text of the ruling.
	Comment string `json:"comment"`
}
