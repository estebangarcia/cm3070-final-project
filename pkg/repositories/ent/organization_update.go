// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/registry"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/user"
)

// OrganizationUpdate is the builder for updating Organization entities.
type OrganizationUpdate struct {
	config
	hooks    []Hook
	mutation *OrganizationMutation
}

// Where appends a list predicates to the OrganizationUpdate builder.
func (ou *OrganizationUpdate) Where(ps ...predicate.Organization) *OrganizationUpdate {
	ou.mutation.Where(ps...)
	return ou
}

// SetName sets the "name" field.
func (ou *OrganizationUpdate) SetName(s string) *OrganizationUpdate {
	ou.mutation.SetName(s)
	return ou
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ou *OrganizationUpdate) SetNillableName(s *string) *OrganizationUpdate {
	if s != nil {
		ou.SetName(*s)
	}
	return ou
}

// SetSlug sets the "slug" field.
func (ou *OrganizationUpdate) SetSlug(s string) *OrganizationUpdate {
	ou.mutation.SetSlug(s)
	return ou
}

// SetNillableSlug sets the "slug" field if the given value is not nil.
func (ou *OrganizationUpdate) SetNillableSlug(s *string) *OrganizationUpdate {
	if s != nil {
		ou.SetSlug(*s)
	}
	return ou
}

// SetIsPersonal sets the "is_personal" field.
func (ou *OrganizationUpdate) SetIsPersonal(b bool) *OrganizationUpdate {
	ou.mutation.SetIsPersonal(b)
	return ou
}

// SetNillableIsPersonal sets the "is_personal" field if the given value is not nil.
func (ou *OrganizationUpdate) SetNillableIsPersonal(b *bool) *OrganizationUpdate {
	if b != nil {
		ou.SetIsPersonal(*b)
	}
	return ou
}

// AddRegistryIDs adds the "registries" edge to the Registry entity by IDs.
func (ou *OrganizationUpdate) AddRegistryIDs(ids ...int) *OrganizationUpdate {
	ou.mutation.AddRegistryIDs(ids...)
	return ou
}

// AddRegistries adds the "registries" edges to the Registry entity.
func (ou *OrganizationUpdate) AddRegistries(r ...*Registry) *OrganizationUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ou.AddRegistryIDs(ids...)
}

// AddMemberIDs adds the "members" edge to the User entity by IDs.
func (ou *OrganizationUpdate) AddMemberIDs(ids ...int) *OrganizationUpdate {
	ou.mutation.AddMemberIDs(ids...)
	return ou
}

// AddMembers adds the "members" edges to the User entity.
func (ou *OrganizationUpdate) AddMembers(u ...*User) *OrganizationUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ou.AddMemberIDs(ids...)
}

// Mutation returns the OrganizationMutation object of the builder.
func (ou *OrganizationUpdate) Mutation() *OrganizationMutation {
	return ou.mutation
}

// ClearRegistries clears all "registries" edges to the Registry entity.
func (ou *OrganizationUpdate) ClearRegistries() *OrganizationUpdate {
	ou.mutation.ClearRegistries()
	return ou
}

// RemoveRegistryIDs removes the "registries" edge to Registry entities by IDs.
func (ou *OrganizationUpdate) RemoveRegistryIDs(ids ...int) *OrganizationUpdate {
	ou.mutation.RemoveRegistryIDs(ids...)
	return ou
}

// RemoveRegistries removes "registries" edges to Registry entities.
func (ou *OrganizationUpdate) RemoveRegistries(r ...*Registry) *OrganizationUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ou.RemoveRegistryIDs(ids...)
}

// ClearMembers clears all "members" edges to the User entity.
func (ou *OrganizationUpdate) ClearMembers() *OrganizationUpdate {
	ou.mutation.ClearMembers()
	return ou
}

// RemoveMemberIDs removes the "members" edge to User entities by IDs.
func (ou *OrganizationUpdate) RemoveMemberIDs(ids ...int) *OrganizationUpdate {
	ou.mutation.RemoveMemberIDs(ids...)
	return ou
}

// RemoveMembers removes "members" edges to User entities.
func (ou *OrganizationUpdate) RemoveMembers(u ...*User) *OrganizationUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ou.RemoveMemberIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ou *OrganizationUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ou.sqlSave, ou.mutation, ou.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ou *OrganizationUpdate) SaveX(ctx context.Context) int {
	affected, err := ou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ou *OrganizationUpdate) Exec(ctx context.Context) error {
	_, err := ou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ou *OrganizationUpdate) ExecX(ctx context.Context) {
	if err := ou.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ou *OrganizationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(organization.Table, organization.Columns, sqlgraph.NewFieldSpec(organization.FieldID, field.TypeInt))
	if ps := ou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ou.mutation.Name(); ok {
		_spec.SetField(organization.FieldName, field.TypeString, value)
	}
	if value, ok := ou.mutation.Slug(); ok {
		_spec.SetField(organization.FieldSlug, field.TypeString, value)
	}
	if value, ok := ou.mutation.IsPersonal(); ok {
		_spec.SetField(organization.FieldIsPersonal, field.TypeBool, value)
	}
	if ou.mutation.RegistriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   organization.RegistriesTable,
			Columns: []string{organization.RegistriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(registry.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.RemovedRegistriesIDs(); len(nodes) > 0 && !ou.mutation.RegistriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   organization.RegistriesTable,
			Columns: []string{organization.RegistriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(registry.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.RegistriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   organization.RegistriesTable,
			Columns: []string{organization.RegistriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(registry.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ou.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   organization.MembersTable,
			Columns: organization.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.RemovedMembersIDs(); len(nodes) > 0 && !ou.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   organization.MembersTable,
			Columns: organization.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   organization.MembersTable,
			Columns: organization.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{organization.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ou.mutation.done = true
	return n, nil
}

// OrganizationUpdateOne is the builder for updating a single Organization entity.
type OrganizationUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *OrganizationMutation
}

// SetName sets the "name" field.
func (ouo *OrganizationUpdateOne) SetName(s string) *OrganizationUpdateOne {
	ouo.mutation.SetName(s)
	return ouo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ouo *OrganizationUpdateOne) SetNillableName(s *string) *OrganizationUpdateOne {
	if s != nil {
		ouo.SetName(*s)
	}
	return ouo
}

// SetSlug sets the "slug" field.
func (ouo *OrganizationUpdateOne) SetSlug(s string) *OrganizationUpdateOne {
	ouo.mutation.SetSlug(s)
	return ouo
}

// SetNillableSlug sets the "slug" field if the given value is not nil.
func (ouo *OrganizationUpdateOne) SetNillableSlug(s *string) *OrganizationUpdateOne {
	if s != nil {
		ouo.SetSlug(*s)
	}
	return ouo
}

// SetIsPersonal sets the "is_personal" field.
func (ouo *OrganizationUpdateOne) SetIsPersonal(b bool) *OrganizationUpdateOne {
	ouo.mutation.SetIsPersonal(b)
	return ouo
}

// SetNillableIsPersonal sets the "is_personal" field if the given value is not nil.
func (ouo *OrganizationUpdateOne) SetNillableIsPersonal(b *bool) *OrganizationUpdateOne {
	if b != nil {
		ouo.SetIsPersonal(*b)
	}
	return ouo
}

// AddRegistryIDs adds the "registries" edge to the Registry entity by IDs.
func (ouo *OrganizationUpdateOne) AddRegistryIDs(ids ...int) *OrganizationUpdateOne {
	ouo.mutation.AddRegistryIDs(ids...)
	return ouo
}

// AddRegistries adds the "registries" edges to the Registry entity.
func (ouo *OrganizationUpdateOne) AddRegistries(r ...*Registry) *OrganizationUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ouo.AddRegistryIDs(ids...)
}

// AddMemberIDs adds the "members" edge to the User entity by IDs.
func (ouo *OrganizationUpdateOne) AddMemberIDs(ids ...int) *OrganizationUpdateOne {
	ouo.mutation.AddMemberIDs(ids...)
	return ouo
}

// AddMembers adds the "members" edges to the User entity.
func (ouo *OrganizationUpdateOne) AddMembers(u ...*User) *OrganizationUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ouo.AddMemberIDs(ids...)
}

// Mutation returns the OrganizationMutation object of the builder.
func (ouo *OrganizationUpdateOne) Mutation() *OrganizationMutation {
	return ouo.mutation
}

// ClearRegistries clears all "registries" edges to the Registry entity.
func (ouo *OrganizationUpdateOne) ClearRegistries() *OrganizationUpdateOne {
	ouo.mutation.ClearRegistries()
	return ouo
}

// RemoveRegistryIDs removes the "registries" edge to Registry entities by IDs.
func (ouo *OrganizationUpdateOne) RemoveRegistryIDs(ids ...int) *OrganizationUpdateOne {
	ouo.mutation.RemoveRegistryIDs(ids...)
	return ouo
}

// RemoveRegistries removes "registries" edges to Registry entities.
func (ouo *OrganizationUpdateOne) RemoveRegistries(r ...*Registry) *OrganizationUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ouo.RemoveRegistryIDs(ids...)
}

// ClearMembers clears all "members" edges to the User entity.
func (ouo *OrganizationUpdateOne) ClearMembers() *OrganizationUpdateOne {
	ouo.mutation.ClearMembers()
	return ouo
}

// RemoveMemberIDs removes the "members" edge to User entities by IDs.
func (ouo *OrganizationUpdateOne) RemoveMemberIDs(ids ...int) *OrganizationUpdateOne {
	ouo.mutation.RemoveMemberIDs(ids...)
	return ouo
}

// RemoveMembers removes "members" edges to User entities.
func (ouo *OrganizationUpdateOne) RemoveMembers(u ...*User) *OrganizationUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ouo.RemoveMemberIDs(ids...)
}

// Where appends a list predicates to the OrganizationUpdate builder.
func (ouo *OrganizationUpdateOne) Where(ps ...predicate.Organization) *OrganizationUpdateOne {
	ouo.mutation.Where(ps...)
	return ouo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ouo *OrganizationUpdateOne) Select(field string, fields ...string) *OrganizationUpdateOne {
	ouo.fields = append([]string{field}, fields...)
	return ouo
}

// Save executes the query and returns the updated Organization entity.
func (ouo *OrganizationUpdateOne) Save(ctx context.Context) (*Organization, error) {
	return withHooks(ctx, ouo.sqlSave, ouo.mutation, ouo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ouo *OrganizationUpdateOne) SaveX(ctx context.Context) *Organization {
	node, err := ouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ouo *OrganizationUpdateOne) Exec(ctx context.Context) error {
	_, err := ouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ouo *OrganizationUpdateOne) ExecX(ctx context.Context) {
	if err := ouo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ouo *OrganizationUpdateOne) sqlSave(ctx context.Context) (_node *Organization, err error) {
	_spec := sqlgraph.NewUpdateSpec(organization.Table, organization.Columns, sqlgraph.NewFieldSpec(organization.FieldID, field.TypeInt))
	id, ok := ouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Organization.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ouo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, organization.FieldID)
		for _, f := range fields {
			if !organization.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != organization.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ouo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ouo.mutation.Name(); ok {
		_spec.SetField(organization.FieldName, field.TypeString, value)
	}
	if value, ok := ouo.mutation.Slug(); ok {
		_spec.SetField(organization.FieldSlug, field.TypeString, value)
	}
	if value, ok := ouo.mutation.IsPersonal(); ok {
		_spec.SetField(organization.FieldIsPersonal, field.TypeBool, value)
	}
	if ouo.mutation.RegistriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   organization.RegistriesTable,
			Columns: []string{organization.RegistriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(registry.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.RemovedRegistriesIDs(); len(nodes) > 0 && !ouo.mutation.RegistriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   organization.RegistriesTable,
			Columns: []string{organization.RegistriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(registry.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.RegistriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   organization.RegistriesTable,
			Columns: []string{organization.RegistriesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(registry.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ouo.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   organization.MembersTable,
			Columns: organization.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.RemovedMembersIDs(); len(nodes) > 0 && !ouo.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   organization.MembersTable,
			Columns: organization.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   organization.MembersTable,
			Columns: organization.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Organization{config: ouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{organization.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ouo.mutation.done = true
	return _node, nil
}
