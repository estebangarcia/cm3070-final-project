// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
)

// Organization is the model entity for the Organization schema.
type Organization struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Slug holds the value of the "slug" field.
	Slug string `json:"slug,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the OrganizationQuery when eager-loading is set.
	Edges        OrganizationEdges `json:"edges"`
	selectValues sql.SelectValues
}

// OrganizationEdges holds the relations/edges for other nodes in the graph.
type OrganizationEdges struct {
	// Registries holds the value of the registries edge.
	Registries []*Registry `json:"registries,omitempty"`
	// Members holds the value of the members edge.
	Members []*User `json:"members,omitempty"`
	// OrgMembers holds the value of the org_members edge.
	OrgMembers []*OrganizationMembership `json:"org_members,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// RegistriesOrErr returns the Registries value or an error if the edge
// was not loaded in eager-loading.
func (e OrganizationEdges) RegistriesOrErr() ([]*Registry, error) {
	if e.loadedTypes[0] {
		return e.Registries, nil
	}
	return nil, &NotLoadedError{edge: "registries"}
}

// MembersOrErr returns the Members value or an error if the edge
// was not loaded in eager-loading.
func (e OrganizationEdges) MembersOrErr() ([]*User, error) {
	if e.loadedTypes[1] {
		return e.Members, nil
	}
	return nil, &NotLoadedError{edge: "members"}
}

// OrgMembersOrErr returns the OrgMembers value or an error if the edge
// was not loaded in eager-loading.
func (e OrganizationEdges) OrgMembersOrErr() ([]*OrganizationMembership, error) {
	if e.loadedTypes[2] {
		return e.OrgMembers, nil
	}
	return nil, &NotLoadedError{edge: "org_members"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Organization) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case organization.FieldID:
			values[i] = new(sql.NullInt64)
		case organization.FieldName, organization.FieldSlug:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Organization fields.
func (o *Organization) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case organization.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			o.ID = int(value.Int64)
		case organization.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				o.Name = value.String
			}
		case organization.FieldSlug:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field slug", values[i])
			} else if value.Valid {
				o.Slug = value.String
			}
		default:
			o.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Organization.
// This includes values selected through modifiers, order, etc.
func (o *Organization) Value(name string) (ent.Value, error) {
	return o.selectValues.Get(name)
}

// QueryRegistries queries the "registries" edge of the Organization entity.
func (o *Organization) QueryRegistries() *RegistryQuery {
	return NewOrganizationClient(o.config).QueryRegistries(o)
}

// QueryMembers queries the "members" edge of the Organization entity.
func (o *Organization) QueryMembers() *UserQuery {
	return NewOrganizationClient(o.config).QueryMembers(o)
}

// QueryOrgMembers queries the "org_members" edge of the Organization entity.
func (o *Organization) QueryOrgMembers() *OrganizationMembershipQuery {
	return NewOrganizationClient(o.config).QueryOrgMembers(o)
}

// Update returns a builder for updating this Organization.
// Note that you need to call Organization.Unwrap() before calling this method if this Organization
// was returned from a transaction, and the transaction was committed or rolled back.
func (o *Organization) Update() *OrganizationUpdateOne {
	return NewOrganizationClient(o.config).UpdateOne(o)
}

// Unwrap unwraps the Organization entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (o *Organization) Unwrap() *Organization {
	_tx, ok := o.config.driver.(*txDriver)
	if !ok {
		panic("ent: Organization is not a transactional entity")
	}
	o.config.driver = _tx.drv
	return o
}

// String implements the fmt.Stringer.
func (o *Organization) String() string {
	var builder strings.Builder
	builder.WriteString("Organization(")
	builder.WriteString(fmt.Sprintf("id=%v, ", o.ID))
	builder.WriteString("name=")
	builder.WriteString(o.Name)
	builder.WriteString(", ")
	builder.WriteString("slug=")
	builder.WriteString(o.Slug)
	builder.WriteByte(')')
	return builder.String()
}

// Organizations is a parsable slice of Organization.
type Organizations []*Organization
