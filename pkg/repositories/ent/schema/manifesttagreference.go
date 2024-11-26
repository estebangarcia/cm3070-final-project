package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ManifestTagReference holds the schema definition for the ManifestTagReference entity.
type ManifestTagReference struct {
	ent.Schema
}

// Fields of the ManifestTagReference.
func (ManifestTagReference) Fields() []ent.Field {
	return []ent.Field{
		field.String("tag"),
	}
}

// Edges of the ManifestTagReference.
func (ManifestTagReference) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("manifests", Manifest.Type).
			Unique(),
	}
}
