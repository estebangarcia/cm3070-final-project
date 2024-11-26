package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Manifest holds the schema definition for the Manifest entity.
type Manifest struct {
	ent.Schema
}

// Fields of the Manifest.
func (Manifest) Fields() []ent.Field {
	return []ent.Field{
		field.String("media_type"),
		field.String("s3_path"),
		field.String("digest"),
	}
}

// Edges of the Manifest.
func (Manifest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tags", ManifestTagReference.Type).
			Ref("manifests"),
		edge.From("repository", Repository.Type).Ref("manifests").Unique(),
	}
}
