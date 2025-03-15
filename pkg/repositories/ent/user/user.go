// Code generated by ent, DO NOT EDIT.

package user

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldGivenName holds the string denoting the given_name field in the database.
	FieldGivenName = "given_name"
	// FieldFamilyName holds the string denoting the family_name field in the database.
	FieldFamilyName = "family_name"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldSub holds the string denoting the sub field in the database.
	FieldSub = "sub"
	// EdgeOrganizations holds the string denoting the organizations edge name in mutations.
	EdgeOrganizations = "organizations"
	// EdgeOrganizationInvites holds the string denoting the organization_invites edge name in mutations.
	EdgeOrganizationInvites = "organization_invites"
	// EdgeJoinedOrganizations holds the string denoting the joined_organizations edge name in mutations.
	EdgeJoinedOrganizations = "joined_organizations"
	// Table holds the table name of the user in the database.
	Table = "users"
	// OrganizationsTable is the table that holds the organizations relation/edge. The primary key declared below.
	OrganizationsTable = "organization_memberships"
	// OrganizationsInverseTable is the table name for the Organization entity.
	// It exists in this package in order to avoid circular dependency with the "organization" package.
	OrganizationsInverseTable = "organizations"
	// OrganizationInvitesTable is the table that holds the organization_invites relation/edge.
	OrganizationInvitesTable = "organization_invites"
	// OrganizationInvitesInverseTable is the table name for the OrganizationInvite entity.
	// It exists in this package in order to avoid circular dependency with the "organizationinvite" package.
	OrganizationInvitesInverseTable = "organization_invites"
	// OrganizationInvitesColumn is the table column denoting the organization_invites relation/edge.
	OrganizationInvitesColumn = "user_id"
	// JoinedOrganizationsTable is the table that holds the joined_organizations relation/edge.
	JoinedOrganizationsTable = "organization_memberships"
	// JoinedOrganizationsInverseTable is the table name for the OrganizationMembership entity.
	// It exists in this package in order to avoid circular dependency with the "organizationmembership" package.
	JoinedOrganizationsInverseTable = "organization_memberships"
	// JoinedOrganizationsColumn is the table column denoting the joined_organizations relation/edge.
	JoinedOrganizationsColumn = "user_id"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldGivenName,
	FieldFamilyName,
	FieldEmail,
	FieldSub,
}

var (
	// OrganizationsPrimaryKey and OrganizationsColumn2 are the table columns denoting the
	// primary key for the organizations relation (M2M).
	OrganizationsPrimaryKey = []string{"user_id", "organization_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByGivenName orders the results by the given_name field.
func ByGivenName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGivenName, opts...).ToFunc()
}

// ByFamilyName orders the results by the family_name field.
func ByFamilyName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFamilyName, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// BySub orders the results by the sub field.
func BySub(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSub, opts...).ToFunc()
}

// ByOrganizationsCount orders the results by organizations count.
func ByOrganizationsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newOrganizationsStep(), opts...)
	}
}

// ByOrganizations orders the results by organizations terms.
func ByOrganizations(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOrganizationsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByOrganizationInvitesCount orders the results by organization_invites count.
func ByOrganizationInvitesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newOrganizationInvitesStep(), opts...)
	}
}

// ByOrganizationInvites orders the results by organization_invites terms.
func ByOrganizationInvites(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOrganizationInvitesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByJoinedOrganizationsCount orders the results by joined_organizations count.
func ByJoinedOrganizationsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newJoinedOrganizationsStep(), opts...)
	}
}

// ByJoinedOrganizations orders the results by joined_organizations terms.
func ByJoinedOrganizations(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newJoinedOrganizationsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newOrganizationsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrganizationsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, OrganizationsTable, OrganizationsPrimaryKey...),
	)
}
func newOrganizationInvitesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrganizationInvitesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, OrganizationInvitesTable, OrganizationInvitesColumn),
	)
}
func newJoinedOrganizationsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(JoinedOrganizationsInverseTable, JoinedOrganizationsColumn),
		sqlgraph.Edge(sqlgraph.O2M, true, JoinedOrganizationsTable, JoinedOrganizationsColumn),
	)
}
