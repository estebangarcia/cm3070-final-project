package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ManifestLayer holds the schema definition for the ManifestLayer entity.
type ManifestLayer struct {
	ent.Schema
}

// Fields of the ManifestLayer.
func (ManifestLayer) Fields() []ent.Field {
	return []ent.Field{
		field.String("media_type"),
		field.String("digest"),
		field.Int32("size"),
		field.JSON("annotations", map[string]string{}),
	}
}

func (ManifestLayer) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("digest").Edges("manifest").Unique(),
	}
}

// Edges of the ManifestLayer.
func (ManifestLayer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("manifest", Manifest.Type).Ref("manifest_layers").Unique(),
	}
}
