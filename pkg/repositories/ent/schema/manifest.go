package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		field.String("artifact_type").Optional(),
		field.String("s3_path"),
		field.String("digest"),
		field.Time("scanned_at").Optional().Nillable(),
		field.Time("uploaded_at").Default(time.Now).Optional().Nillable(),
	}
}

// Edges of the Manifest.
func (Manifest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tags", ManifestTagReference.Type).
			Ref("manifests"),
		edge.From("repository", Repository.Type).Ref("manifests").Unique(),
		edge.To("subject", Manifest.Type),
		edge.From("referer", Manifest.Type).Ref("subject"),
		edge.To("manifest_layers", ManifestLayer.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("vulnerabilities", Vulnerability.Type).Ref("manifests"),
	}
}
