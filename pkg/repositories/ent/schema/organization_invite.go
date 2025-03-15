package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/oklog/ulid/v2"
)

// OrganizationInvite holds the schema definition for the OrganizationInvite entity.
type OrganizationInvite struct {
	ent.Schema
}

// Fields of the OrganizationInvite.
func (OrganizationInvite) Fields() []ent.Field {
	return []ent.Field{
		field.String("invite_id").DefaultFunc(
			func() string {
				return ulid.Make().String()
			},
		),
		field.Int("organization_id"),
		field.Int("user_id").Nillable().Optional(),
		field.String("email").Nillable().Optional(),
		field.Enum("role").Values(RoleNames...),
	}
}

func (OrganizationInvite) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("invite_id").Unique(),
		index.Fields("user_id", "organization_id").Unique(),
		index.Fields("email", "organization_id").Unique(),
	}
}

// Edges of the OrganizationInvite.
func (OrganizationInvite) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).Ref("organization_invites").Field("organization_id").Required().Unique(),
		edge.From("invitee", User.Type).Ref("organization_invites").Field("user_id").Unique(),
	}
}
