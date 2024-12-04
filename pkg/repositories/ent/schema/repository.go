package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Manifest holds the schema definition for the Repository entity.
type Repository struct {
	ent.Schema
}

// Fields of the Repository.
func (Repository) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

func (Repository) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Edges("registry").Unique(),
	}
}

// Edges of the Repository.
func (Repository) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("manifests", Manifest.Type),
		edge.From("registry", Registry.Type).Ref("repositories").Unique(),
	}
}
