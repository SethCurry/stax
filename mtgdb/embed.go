package mtgdb

// only used to embed.
import _ "embed"

//go:embed schema.sql
var schema []byte
