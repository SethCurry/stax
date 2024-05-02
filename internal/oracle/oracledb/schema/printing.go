package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Printing struct {
	ent.Schema
}

func (Printing) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("rarity").Values("common", "uncommon", "rare", "mythic", "special", "bonus"),
	}
}

func (Printing) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("artist", Artist.Type).Unique(),
		edge.To("set", Set.Type).Unique(),
		edge.To("card_face", CardFace.Type).Unique(),
		edge.From("images", PrintingImage.Type).Ref("printing"),
	}
}
