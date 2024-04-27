# stax

[![Go Test](https://github.com/SethCurry/stax/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/SethCurry/stax/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/SethCurry/stax/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/SethCurry/stax/actions/workflows/lint.yml)

`stax` is a Swiss Army Knife for Magic: The Gathering.

Once this project is ready for use, you will be able to download
it from the releases page of this repo.

## Examples

### Scryfall

Searching for a card:

```bash
# Find all cards that have the type creature and a CMC less than 5
stax scryfall card search "cmc<5 AND type:creature"
```
