// Code generated by ent, DO NOT EDIT.

package organization

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the organization type in the database.
	Label = "organization"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldSlug holds the string denoting the slug field in the database.
	FieldSlug = "slug"
	// FieldIsPersonal holds the string denoting the is_personal field in the database.
	FieldIsPersonal = "is_personal"
	// EdgeRegistries holds the string denoting the registries edge name in mutations.
	EdgeRegistries = "registries"
	// EdgeMembers holds the string denoting the members edge name in mutations.
	EdgeMembers = "members"
	// EdgeOrganizationInvites holds the string denoting the organization_invites edge name in mutations.
	EdgeOrganizationInvites = "organization_invites"
	// EdgeOrgMembers holds the string denoting the org_members edge name in mutations.
	EdgeOrgMembers = "org_members"
	// Table holds the table name of the organization in the database.
	Table = "organizations"
	// RegistriesTable is the table that holds the registries relation/edge.
	RegistriesTable = "registries"
	// RegistriesInverseTable is the table name for the Registry entity.
	// It exists in this package in order to avoid circular dependency with the "registry" package.
	RegistriesInverseTable = "registries"
	// RegistriesColumn is the table column denoting the registries relation/edge.
	RegistriesColumn = "organization_registries"
	// MembersTable is the table that holds the members relation/edge. The primary key declared below.
	MembersTable = "organization_memberships"
	// MembersInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	MembersInverseTable = "users"
	// OrganizationInvitesTable is the table that holds the organization_invites relation/edge.
	OrganizationInvitesTable = "organization_invites"
	// OrganizationInvitesInverseTable is the table name for the OrganizationInvite entity.
	// It exists in this package in order to avoid circular dependency with the "organizationinvite" package.
	OrganizationInvitesInverseTable = "organization_invites"
	// OrganizationInvitesColumn is the table column denoting the organization_invites relation/edge.
	OrganizationInvitesColumn = "organization_id"
	// OrgMembersTable is the table that holds the org_members relation/edge.
	OrgMembersTable = "organization_memberships"
	// OrgMembersInverseTable is the table name for the OrganizationMembership entity.
	// It exists in this package in order to avoid circular dependency with the "organizationmembership" package.
	OrgMembersInverseTable = "organization_memberships"
	// OrgMembersColumn is the table column denoting the org_members relation/edge.
	OrgMembersColumn = "organization_id"
)

// Columns holds all SQL columns for organization fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldSlug,
	FieldIsPersonal,
}

var (
	// MembersPrimaryKey and MembersColumn2 are the table columns denoting the
	// primary key for the members relation (M2M).
	MembersPrimaryKey = []string{"user_id", "organization_id"}
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

// OrderOption defines the ordering options for the Organization queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// BySlug orders the results by the slug field.
func BySlug(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSlug, opts...).ToFunc()
}

// ByIsPersonal orders the results by the is_personal field.
func ByIsPersonal(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsPersonal, opts...).ToFunc()
}

// ByRegistriesCount orders the results by registries count.
func ByRegistriesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newRegistriesStep(), opts...)
	}
}

// ByRegistries orders the results by registries terms.
func ByRegistries(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRegistriesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByMembersCount orders the results by members count.
func ByMembersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMembersStep(), opts...)
	}
}

// ByMembers orders the results by members terms.
func ByMembers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMembersStep(), append([]sql.OrderTerm{term}, terms...)...)
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

// ByOrgMembersCount orders the results by org_members count.
func ByOrgMembersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newOrgMembersStep(), opts...)
	}
}

// ByOrgMembers orders the results by org_members terms.
func ByOrgMembers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOrgMembersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newRegistriesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RegistriesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, RegistriesTable, RegistriesColumn),
	)
}
func newMembersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MembersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, MembersTable, MembersPrimaryKey...),
	)
}
func newOrganizationInvitesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrganizationInvitesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, OrganizationInvitesTable, OrganizationInvitesColumn),
	)
}
func newOrgMembersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrgMembersInverseTable, OrgMembersColumn),
		sqlgraph.Edge(sqlgraph.O2M, true, OrgMembersTable, OrgMembersColumn),
	)
}
