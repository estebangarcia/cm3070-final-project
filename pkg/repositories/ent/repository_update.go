// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
)

// RepositoryUpdate is the builder for updating Repository entities.
type RepositoryUpdate struct {
	config
	hooks    []Hook
	mutation *RepositoryMutation
}

// Where appends a list predicates to the RepositoryUpdate builder.
func (ru *RepositoryUpdate) Where(ps ...predicate.Repository) *RepositoryUpdate {
	ru.mutation.Where(ps...)
	return ru
}

// SetName sets the "name" field.
func (ru *RepositoryUpdate) SetName(s string) *RepositoryUpdate {
	ru.mutation.SetName(s)
	return ru
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ru *RepositoryUpdate) SetNillableName(s *string) *RepositoryUpdate {
	if s != nil {
		ru.SetName(*s)
	}
	return ru
}

// AddManifestIDs adds the "manifests" edge to the Manifest entity by IDs.
func (ru *RepositoryUpdate) AddManifestIDs(ids ...int) *RepositoryUpdate {
	ru.mutation.AddManifestIDs(ids...)
	return ru
}

// AddManifests adds the "manifests" edges to the Manifest entity.
func (ru *RepositoryUpdate) AddManifests(m ...*Manifest) *RepositoryUpdate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return ru.AddManifestIDs(ids...)
}

// Mutation returns the RepositoryMutation object of the builder.
func (ru *RepositoryUpdate) Mutation() *RepositoryMutation {
	return ru.mutation
}

// ClearManifests clears all "manifests" edges to the Manifest entity.
func (ru *RepositoryUpdate) ClearManifests() *RepositoryUpdate {
	ru.mutation.ClearManifests()
	return ru
}

// RemoveManifestIDs removes the "manifests" edge to Manifest entities by IDs.
func (ru *RepositoryUpdate) RemoveManifestIDs(ids ...int) *RepositoryUpdate {
	ru.mutation.RemoveManifestIDs(ids...)
	return ru
}

// RemoveManifests removes "manifests" edges to Manifest entities.
func (ru *RepositoryUpdate) RemoveManifests(m ...*Manifest) *RepositoryUpdate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return ru.RemoveManifestIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ru *RepositoryUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ru.sqlSave, ru.mutation, ru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ru *RepositoryUpdate) SaveX(ctx context.Context) int {
	affected, err := ru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ru *RepositoryUpdate) Exec(ctx context.Context) error {
	_, err := ru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ru *RepositoryUpdate) ExecX(ctx context.Context) {
	if err := ru.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ru *RepositoryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(repository.Table, repository.Columns, sqlgraph.NewFieldSpec(repository.FieldID, field.TypeInt))
	if ps := ru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.Name(); ok {
		_spec.SetField(repository.FieldName, field.TypeString, value)
	}
	if ru.mutation.ManifestsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repository.ManifestsTable,
			Columns: []string{repository.ManifestsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.RemovedManifestsIDs(); len(nodes) > 0 && !ru.mutation.ManifestsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repository.ManifestsTable,
			Columns: []string{repository.ManifestsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.ManifestsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repository.ManifestsTable,
			Columns: []string{repository.ManifestsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{repository.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ru.mutation.done = true
	return n, nil
}

// RepositoryUpdateOne is the builder for updating a single Repository entity.
type RepositoryUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *RepositoryMutation
}

// SetName sets the "name" field.
func (ruo *RepositoryUpdateOne) SetName(s string) *RepositoryUpdateOne {
	ruo.mutation.SetName(s)
	return ruo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ruo *RepositoryUpdateOne) SetNillableName(s *string) *RepositoryUpdateOne {
	if s != nil {
		ruo.SetName(*s)
	}
	return ruo
}

// AddManifestIDs adds the "manifests" edge to the Manifest entity by IDs.
func (ruo *RepositoryUpdateOne) AddManifestIDs(ids ...int) *RepositoryUpdateOne {
	ruo.mutation.AddManifestIDs(ids...)
	return ruo
}

// AddManifests adds the "manifests" edges to the Manifest entity.
func (ruo *RepositoryUpdateOne) AddManifests(m ...*Manifest) *RepositoryUpdateOne {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return ruo.AddManifestIDs(ids...)
}

// Mutation returns the RepositoryMutation object of the builder.
func (ruo *RepositoryUpdateOne) Mutation() *RepositoryMutation {
	return ruo.mutation
}

// ClearManifests clears all "manifests" edges to the Manifest entity.
func (ruo *RepositoryUpdateOne) ClearManifests() *RepositoryUpdateOne {
	ruo.mutation.ClearManifests()
	return ruo
}

// RemoveManifestIDs removes the "manifests" edge to Manifest entities by IDs.
func (ruo *RepositoryUpdateOne) RemoveManifestIDs(ids ...int) *RepositoryUpdateOne {
	ruo.mutation.RemoveManifestIDs(ids...)
	return ruo
}

// RemoveManifests removes "manifests" edges to Manifest entities.
func (ruo *RepositoryUpdateOne) RemoveManifests(m ...*Manifest) *RepositoryUpdateOne {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return ruo.RemoveManifestIDs(ids...)
}

// Where appends a list predicates to the RepositoryUpdate builder.
func (ruo *RepositoryUpdateOne) Where(ps ...predicate.Repository) *RepositoryUpdateOne {
	ruo.mutation.Where(ps...)
	return ruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ruo *RepositoryUpdateOne) Select(field string, fields ...string) *RepositoryUpdateOne {
	ruo.fields = append([]string{field}, fields...)
	return ruo
}

// Save executes the query and returns the updated Repository entity.
func (ruo *RepositoryUpdateOne) Save(ctx context.Context) (*Repository, error) {
	return withHooks(ctx, ruo.sqlSave, ruo.mutation, ruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ruo *RepositoryUpdateOne) SaveX(ctx context.Context) *Repository {
	node, err := ruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ruo *RepositoryUpdateOne) Exec(ctx context.Context) error {
	_, err := ruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *RepositoryUpdateOne) ExecX(ctx context.Context) {
	if err := ruo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ruo *RepositoryUpdateOne) sqlSave(ctx context.Context) (_node *Repository, err error) {
	_spec := sqlgraph.NewUpdateSpec(repository.Table, repository.Columns, sqlgraph.NewFieldSpec(repository.FieldID, field.TypeInt))
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Repository.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, repository.FieldID)
		for _, f := range fields {
			if !repository.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != repository.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruo.mutation.Name(); ok {
		_spec.SetField(repository.FieldName, field.TypeString, value)
	}
	if ruo.mutation.ManifestsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repository.ManifestsTable,
			Columns: []string{repository.ManifestsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.RemovedManifestsIDs(); len(nodes) > 0 && !ruo.mutation.ManifestsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repository.ManifestsTable,
			Columns: []string{repository.ManifestsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.ManifestsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   repository.ManifestsTable,
			Columns: []string{repository.ManifestsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Repository{config: ruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{repository.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ruo.mutation.done = true
	return _node, nil
}
