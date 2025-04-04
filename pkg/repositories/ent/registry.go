// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/registry"
)

// Registry is the model entity for the Registry schema.
type Registry struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Slug holds the value of the "slug" field.
	Slug string `json:"slug,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the RegistryQuery when eager-loading is set.
	Edges                   RegistryEdges `json:"edges"`
	organization_registries *int
	selectValues            sql.SelectValues
}

// RegistryEdges holds the relations/edges for other nodes in the graph.
type RegistryEdges struct {
	// Repositories holds the value of the repositories edge.
	Repositories []*Repository `json:"repositories,omitempty"`
	// Organization holds the value of the organization edge.
	Organization *Organization `json:"organization,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// RepositoriesOrErr returns the Repositories value or an error if the edge
// was not loaded in eager-loading.
func (e RegistryEdges) RepositoriesOrErr() ([]*Repository, error) {
	if e.loadedTypes[0] {
		return e.Repositories, nil
	}
	return nil, &NotLoadedError{edge: "repositories"}
}

// OrganizationOrErr returns the Organization value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RegistryEdges) OrganizationOrErr() (*Organization, error) {
	if e.Organization != nil {
		return e.Organization, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: organization.Label}
	}
	return nil, &NotLoadedError{edge: "organization"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Registry) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case registry.FieldID:
			values[i] = new(sql.NullInt64)
		case registry.FieldName, registry.FieldSlug:
			values[i] = new(sql.NullString)
		case registry.ForeignKeys[0]: // organization_registries
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Registry fields.
func (r *Registry) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case registry.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			r.ID = int(value.Int64)
		case registry.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				r.Name = value.String
			}
		case registry.FieldSlug:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field slug", values[i])
			} else if value.Valid {
				r.Slug = value.String
			}
		case registry.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field organization_registries", value)
			} else if value.Valid {
				r.organization_registries = new(int)
				*r.organization_registries = int(value.Int64)
			}
		default:
			r.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Registry.
// This includes values selected through modifiers, order, etc.
func (r *Registry) Value(name string) (ent.Value, error) {
	return r.selectValues.Get(name)
}

// QueryRepositories queries the "repositories" edge of the Registry entity.
func (r *Registry) QueryRepositories() *RepositoryQuery {
	return NewRegistryClient(r.config).QueryRepositories(r)
}

// QueryOrganization queries the "organization" edge of the Registry entity.
func (r *Registry) QueryOrganization() *OrganizationQuery {
	return NewRegistryClient(r.config).QueryOrganization(r)
}

// Update returns a builder for updating this Registry.
// Note that you need to call Registry.Unwrap() before calling this method if this Registry
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *Registry) Update() *RegistryUpdateOne {
	return NewRegistryClient(r.config).UpdateOne(r)
}

// Unwrap unwraps the Registry entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *Registry) Unwrap() *Registry {
	_tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: Registry is not a transactional entity")
	}
	r.config.driver = _tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *Registry) String() string {
	var builder strings.Builder
	builder.WriteString("Registry(")
	builder.WriteString(fmt.Sprintf("id=%v, ", r.ID))
	builder.WriteString("name=")
	builder.WriteString(r.Name)
	builder.WriteString(", ")
	builder.WriteString("slug=")
	builder.WriteString(r.Slug)
	builder.WriteByte(')')
	return builder.String()
}

// Registries is a parsable slice of Registry.
type Registries []*Registry
