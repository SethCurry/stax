package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

const (
	ArtistNameMinLen = 3
	ArtistNameMaxLen = 255
)

type Artist struct {
	ent.Schema
}

func (Artist) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MinLen(ArtistNameMinLen).MaxLen(ArtistNameMaxLen),
	}
}

func (Artist) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("printings", Printing.Type).Ref("artist"),
	}
}
