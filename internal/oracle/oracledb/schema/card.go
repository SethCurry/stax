package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

const (
	CardNameMinLen = 3
	CardNameMaxLen = 255
)

type Card struct {
	ent.Schema
}

func (Card) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MinLen(CardNameMinLen).MaxLen(CardNameMaxLen),
		field.String("oracle_id").NotEmpty(),
		field.Uint8("color_identity"),
	}
}

func (Card) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("faces", CardFace.Type).Ref("card"),
		edge.From("rulings", Ruling.Type).Ref("card"),
	}
}
