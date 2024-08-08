# stax

[![Go Test](https://github.com/SethCurry/stax/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/SethCurry/stax/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/SethCurry/stax/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/SethCurry/stax/actions/workflows/lint.yml)

`stax` is a Swiss Army Knife for Magic: The Gathering.

Once this project is ready for use, you will be able to download
it from the releases page of this repo.

## Features

- Scryfall API Client
  - Search for cards
  - Get rulings for cards
- Self-hosted Scryfall
  - Can import Scryfall's bulk data exports
  - Exposes the same REST API as scryfall.com

## Examples

### Scryfall

`stax` has built-in support for the [Scryfall API](https://scryfall.com/docs/api).
You can use it to search for cards, download bulk data, and more.

#### Searching for a Card

When searching for a card, the query is passed to the Scryfall API as is.
The best resource for learning how to write these queries is the [Scryfall Search Reference](https://scryfall.com/docs/syntax).

```bash
# Find all cards that have the type creature and a CMC less than 5
# You don't have to use quotes, but your shell may interpret symbols like < and >
# as IO redirection.
stax scryfall search "cmc<5 AND type:creature"
```

#### Getting Rulings For A Card

You can retrieve the rulings for a card by name using the `rulings` subcommand.

```bash
# Print all the rulings for the card Winter Orb
stax scryfall rulings Winter Orb
```

### Generating HTML From The Comprehensive Rules

To generate HTML from the comprehensive rules, you first need to download a copy of the rules from [Wizards of the Coast's website](https://magic.wizards.com/en/rules).
Click on the `TXT` link to download the rules as a text file.

Then run:

```bash
stax rules html /path/to/rules.txt
```

You can specify a path to output the HTML file to using the `--output` flag.

```bash
stax rules html /path/to/rules.txt --output /path/to/output.html
```

### Self-hosted Scryfall API

To start a self-hosted Scryfall API, you first need to download and load the bulk data.

This command will download the latest bulk data export from Scryfall and load it into a SQLite database.

```bash
stax bones load --http

# or if you want to manually download a bulk data file

stax bones load /path/to/bulk/data/file.json
```

Once the bulk data is loaded, you can start the API by running:

```bash
stax api
```

If you want to run it on a different port or IP (default is 0.0.0.0:8765), you can use the `--listen` flag:

```bash
stax api --listen 127.0.0.1:8766
```
