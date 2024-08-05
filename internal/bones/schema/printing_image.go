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
		field.Enum("image_type").Values("small", "normal", "large", "png", "art_crop", "border_crop"),
		field.String("local_path").Optional(),
	}
}

func (PrintingImage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("printing", Printing.Type).Unique(),
	}
}
