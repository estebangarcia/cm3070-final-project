package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrganizationMembership holds the schema definition for the OrganizationMembership entity.
type OrganizationMembership struct {
	ent.Schema
}

func (OrganizationMembership) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("user_id", "organization_id"),
	}
}

// Fields of the OrganizationMembership.
func (OrganizationMembership) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("organization_id"),
		field.Int("role"),
	}
}

// Edges of the OrganizationMembership.
func (OrganizationMembership) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id"),
		edge.To("organization", Organization.Type).
			Required().
			Unique().
			Field("organization_id"),
	}
}
