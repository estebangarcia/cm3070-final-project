package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
)

// Misconfiguration holds the schema definition for the Misconfiguration entity.
type Misconfiguration struct {
	ent.Schema
}

// Fields of the Misconfiguration.
func (Misconfiguration) Fields() []ent.Field {
	return []ent.Field{
		field.String("misconfiguration_id").Unique(),
		field.String("misconfiguration_url_details"),
		field.String("title"),
		field.Enum("severity").Values(dbTypes.SeverityNames...),
	}
}

// Edges of the Misconfiguration.
func (Misconfiguration) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("manifest_misconfigurations", ManifestMisconfiguration.Type),
	}
}
