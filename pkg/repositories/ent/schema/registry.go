package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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

// Edges of the Registry.
func (Registry) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).Ref("registries").Unique(),
	}
}
