// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifestlayer"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
)

// ManifestLayerDelete is the builder for deleting a ManifestLayer entity.
type ManifestLayerDelete struct {
	config
	hooks    []Hook
	mutation *ManifestLayerMutation
}

// Where appends a list predicates to the ManifestLayerDelete builder.
func (mld *ManifestLayerDelete) Where(ps ...predicate.ManifestLayer) *ManifestLayerDelete {
	mld.mutation.Where(ps...)
	return mld
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mld *ManifestLayerDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, mld.sqlExec, mld.mutation, mld.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (mld *ManifestLayerDelete) ExecX(ctx context.Context) int {
	n, err := mld.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mld *ManifestLayerDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(manifestlayer.Table, sqlgraph.NewFieldSpec(manifestlayer.FieldID, field.TypeInt))
	if ps := mld.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, mld.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	mld.mutation.done = true
	return affected, err
}

// ManifestLayerDeleteOne is the builder for deleting a single ManifestLayer entity.
type ManifestLayerDeleteOne struct {
	mld *ManifestLayerDelete
}

// Where appends a list predicates to the ManifestLayerDelete builder.
func (mldo *ManifestLayerDeleteOne) Where(ps ...predicate.ManifestLayer) *ManifestLayerDeleteOne {
	mldo.mld.mutation.Where(ps...)
	return mldo
}

// Exec executes the deletion query.
func (mldo *ManifestLayerDeleteOne) Exec(ctx context.Context) error {
	n, err := mldo.mld.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{manifestlayer.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mldo *ManifestLayerDeleteOne) ExecX(ctx context.Context) {
	if err := mldo.Exec(ctx); err != nil {
		panic(err)
	}
}
