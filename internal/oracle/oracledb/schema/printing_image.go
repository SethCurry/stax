package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type PrintingImage struct {
	ent.Schema
}

func (PrintingImage) Fields() []ent.Field {
	return []ent.Field{
		field.String("url").NotEmpty(),
	}
}

func (PrintingImage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("printing", Printing.Type).Unique(),
	}
}
