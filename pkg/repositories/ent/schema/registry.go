package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Registry holds the schema definition for the Registry entity.
type Registry struct {
	ent.Schema
}

// Fields of the Registry.
func (Registry) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug"),
	}
}

func (Registry) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug").Edges("organization").Unique(),
	}
}

// Edges of the Registry.
func (Registry) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("repositories", Repository.Type),
		edge.From("organization", Organization.Type).Ref("registries").Unique(),
	}
}
