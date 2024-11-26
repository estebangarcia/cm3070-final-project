// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifesttagreference"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
)

// ManifestTagReferenceDelete is the builder for deleting a ManifestTagReference entity.
type ManifestTagReferenceDelete struct {
	config
	hooks    []Hook
	mutation *ManifestTagReferenceMutation
}

// Where appends a list predicates to the ManifestTagReferenceDelete builder.
func (mtrd *ManifestTagReferenceDelete) Where(ps ...predicate.ManifestTagReference) *ManifestTagReferenceDelete {
	mtrd.mutation.Where(ps...)
	return mtrd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mtrd *ManifestTagReferenceDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, mtrd.sqlExec, mtrd.mutation, mtrd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (mtrd *ManifestTagReferenceDelete) ExecX(ctx context.Context) int {
	n, err := mtrd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mtrd *ManifestTagReferenceDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(manifesttagreference.Table, sqlgraph.NewFieldSpec(manifesttagreference.FieldID, field.TypeInt))
	if ps := mtrd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, mtrd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	mtrd.mutation.done = true
	return affected, err
}

// ManifestTagReferenceDeleteOne is the builder for deleting a single ManifestTagReference entity.
type ManifestTagReferenceDeleteOne struct {
	mtrd *ManifestTagReferenceDelete
}

// Where appends a list predicates to the ManifestTagReferenceDelete builder.
func (mtrdo *ManifestTagReferenceDeleteOne) Where(ps ...predicate.ManifestTagReference) *ManifestTagReferenceDeleteOne {
	mtrdo.mtrd.mutation.Where(ps...)
	return mtrdo
}

// Exec executes the deletion query.
func (mtrdo *ManifestTagReferenceDeleteOne) Exec(ctx context.Context) error {
	n, err := mtrdo.mtrd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{manifesttagreference.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mtrdo *ManifestTagReferenceDeleteOne) ExecX(ctx context.Context) {
	if err := mtrdo.Exec(ctx); err != nil {
		panic(err)
	}
}
