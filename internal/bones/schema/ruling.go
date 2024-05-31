package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Ruling struct {
	ent.Schema
}

func (Ruling) Fields() []ent.Field {
	return []ent.Field{
		field.String("text").NotEmpty(),
		field.Time("date"),
	}
}

func (Ruling) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("card", Card.Type).Unique(),
	}
}
