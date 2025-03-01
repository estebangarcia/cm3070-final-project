package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Misconfiguration holds the schema definition for the Misconfiguration entity.
type ManifestMisconfiguration struct {
	ent.Schema
}

// Fields of the Misconfiguration.
func (ManifestMisconfiguration) Fields() []ent.Field {
	return []ent.Field{
		field.String("target_file"),
		field.String("message"),
		field.String("resolution"),
		field.Int("manifest_id"),
		field.Int("misconfiguration_id").Optional(),
	}
}

func (ManifestMisconfiguration) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("manifest_id", "misconfiguration_id", "target_file").Unique(),
	}
}

// Edges of the Misconfiguration.
func (ManifestMisconfiguration) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("misconfiguration", Misconfiguration.Type).Ref("manifest_misconfigurations").Unique().Field("misconfiguration_id"),
	}
}
