# stax

[![Go Test](https://github.com/SethCurry/stax/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/SethCurry/stax/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/SethCurry/stax/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/SethCurry/stax/actions/workflows/lint.yml)

`stax` is a Swiss Army Knife for Magic: The Gathering.

Once this project is ready for use, you will be able to download
it from the releases page of this repo.

## Examples

### Scryfall

`stax` has built-in support for the [Scryfall API](https://scryfall.com/docs/api).
You can use it to search for cards, download bulk data, and more.

#### Searching for a Card

When searching for a card, the query is passed to the Scryfall API as is.
The best resource for learning how to write these queries is the [Scryfall Search Reference](https://scryfall.com/docs/syntax).

```bash
# Find all cards that have the type creature and a CMC less than 5
stax scryfall card search "cmc<5 AND type:creature"
```

#### Getting Rulings For A Card
