package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

const (
	CardFaceNameMinLen = 1
	CardFaceNameMaxLen = 255
)

type CardFace struct {
	ent.Schema
}

func (CardFace) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MinLen(CardFaceNameMinLen).MaxLen(CardFaceNameMaxLen),
		field.String("flavor_text").MaxLen(1000),
		field.String("oracle_text").MaxLen(1000),
		field.String("language"),
		field.Float32("cmc"),
		field.String("power"),
		field.String("toughness"),
		field.String("loyalty"),
		field.String("mana_cost"),
		field.String("type_line"),
		field.String("colors"),
	}
}

func (CardFace) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("card", Card.Type).Unique(),
		edge.From("printings", Printing.Type).Ref("card_face"),
	}
}
