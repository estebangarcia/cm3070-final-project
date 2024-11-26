// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
)

// ManifestDelete is the builder for deleting a Manifest entity.
type ManifestDelete struct {
	config
	hooks    []Hook
	mutation *ManifestMutation
}

// Where appends a list predicates to the ManifestDelete builder.
func (md *ManifestDelete) Where(ps ...predicate.Manifest) *ManifestDelete {
	md.mutation.Where(ps...)
	return md
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (md *ManifestDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, md.sqlExec, md.mutation, md.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (md *ManifestDelete) ExecX(ctx context.Context) int {
	n, err := md.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (md *ManifestDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(manifest.Table, sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt))
	if ps := md.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, md.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	md.mutation.done = true
	return affected, err
}

// ManifestDeleteOne is the builder for deleting a single Manifest entity.
type ManifestDeleteOne struct {
	md *ManifestDelete
}

// Where appends a list predicates to the ManifestDelete builder.
func (mdo *ManifestDeleteOne) Where(ps ...predicate.Manifest) *ManifestDeleteOne {
	mdo.md.mutation.Where(ps...)
	return mdo
}

// Exec executes the deletion query.
func (mdo *ManifestDeleteOne) Exec(ctx context.Context) error {
	n, err := mdo.md.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{manifest.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mdo *ManifestDeleteOne) ExecX(ctx context.Context) {
	if err := mdo.Exec(ctx); err != nil {
		panic(err)
	}
}
