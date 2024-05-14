package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

const (
	SetNameMinLen = 3
	SetNameMaxLen = 255
	SetCodeMinLen = 3
	SetCodeMaxLen = 255
)

type Set struct {
	ent.Schema
}

func (Set) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MinLen(SetNameMinLen).MaxLen(SetNameMaxLen),
		field.String("code").NotEmpty().MinLen(SetCodeMinLen).MaxLen(SetCodeMaxLen),
	}
}

func (Set) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("printings", Printing.Type).Ref("set"),
	}
}
