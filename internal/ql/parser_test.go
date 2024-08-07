package ql

func ExampleParseQuery() {
	query := "name=\"Static Orb\" AND cmc<4"

	parsed, err := ParseQuery(query)
	if err != nil {
		panic(err)
	}

	parsed.Predicate()

	// parsed.Predicate() can then be used to filter a card query,
	// like so:
	// cards, err := dbClient.Cards().Query().Where(parsed.Predicate()).All(ctx)
}
