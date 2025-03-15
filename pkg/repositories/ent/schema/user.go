package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("given_name"),
		field.String("family_name"),
		field.String("email"),
		field.String("sub"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("sub").Unique(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("organizations", Organization.Type).
			Through("joined_organizations", OrganizationMembership.Type),
		edge.To("organization_invites", OrganizationInvite.Type),
	}
}
