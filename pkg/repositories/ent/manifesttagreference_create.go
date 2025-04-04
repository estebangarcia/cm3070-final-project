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
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifesttagreference"
)

// ManifestTagReferenceCreate is the builder for creating a ManifestTagReference entity.
type ManifestTagReferenceCreate struct {
	config
	mutation *ManifestTagReferenceMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetTag sets the "tag" field.
func (mtrc *ManifestTagReferenceCreate) SetTag(s string) *ManifestTagReferenceCreate {
	mtrc.mutation.SetTag(s)
	return mtrc
}

// SetManifestsID sets the "manifests" edge to the Manifest entity by ID.
func (mtrc *ManifestTagReferenceCreate) SetManifestsID(id int) *ManifestTagReferenceCreate {
	mtrc.mutation.SetManifestsID(id)
	return mtrc
}

// SetNillableManifestsID sets the "manifests" edge to the Manifest entity by ID if the given value is not nil.
func (mtrc *ManifestTagReferenceCreate) SetNillableManifestsID(id *int) *ManifestTagReferenceCreate {
	if id != nil {
		mtrc = mtrc.SetManifestsID(*id)
	}
	return mtrc
}

// SetManifests sets the "manifests" edge to the Manifest entity.
func (mtrc *ManifestTagReferenceCreate) SetManifests(m *Manifest) *ManifestTagReferenceCreate {
	return mtrc.SetManifestsID(m.ID)
}

// Mutation returns the ManifestTagReferenceMutation object of the builder.
func (mtrc *ManifestTagReferenceCreate) Mutation() *ManifestTagReferenceMutation {
	return mtrc.mutation
}

// Save creates the ManifestTagReference in the database.
func (mtrc *ManifestTagReferenceCreate) Save(ctx context.Context) (*ManifestTagReference, error) {
	return withHooks(ctx, mtrc.sqlSave, mtrc.mutation, mtrc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mtrc *ManifestTagReferenceCreate) SaveX(ctx context.Context) *ManifestTagReference {
	v, err := mtrc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mtrc *ManifestTagReferenceCreate) Exec(ctx context.Context) error {
	_, err := mtrc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mtrc *ManifestTagReferenceCreate) ExecX(ctx context.Context) {
	if err := mtrc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mtrc *ManifestTagReferenceCreate) check() error {
	if _, ok := mtrc.mutation.Tag(); !ok {
		return &ValidationError{Name: "tag", err: errors.New(`ent: missing required field "ManifestTagReference.tag"`)}
	}
	return nil
}

func (mtrc *ManifestTagReferenceCreate) sqlSave(ctx context.Context) (*ManifestTagReference, error) {
	if err := mtrc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mtrc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mtrc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	mtrc.mutation.id = &_node.ID
	mtrc.mutation.done = true
	return _node, nil
}

func (mtrc *ManifestTagReferenceCreate) createSpec() (*ManifestTagReference, *sqlgraph.CreateSpec) {
	var (
		_node = &ManifestTagReference{config: mtrc.config}
		_spec = sqlgraph.NewCreateSpec(manifesttagreference.Table, sqlgraph.NewFieldSpec(manifesttagreference.FieldID, field.TypeInt))
	)
	_spec.OnConflict = mtrc.conflict
	if value, ok := mtrc.mutation.Tag(); ok {
		_spec.SetField(manifesttagreference.FieldTag, field.TypeString, value)
		_node.Tag = value
	}
	if nodes := mtrc.mutation.ManifestsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   manifesttagreference.ManifestsTable,
			Columns: []string{manifesttagreference.ManifestsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.manifest_tag_reference_manifests = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ManifestTagReference.Create().
//		SetTag(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ManifestTagReferenceUpsert) {
//			SetTag(v+v).
//		}).
//		Exec(ctx)
func (mtrc *ManifestTagReferenceCreate) OnConflict(opts ...sql.ConflictOption) *ManifestTagReferenceUpsertOne {
	mtrc.conflict = opts
	return &ManifestTagReferenceUpsertOne{
		create: mtrc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ManifestTagReference.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mtrc *ManifestTagReferenceCreate) OnConflictColumns(columns ...string) *ManifestTagReferenceUpsertOne {
	mtrc.conflict = append(mtrc.conflict, sql.ConflictColumns(columns...))
	return &ManifestTagReferenceUpsertOne{
		create: mtrc,
	}
}

type (
	// ManifestTagReferenceUpsertOne is the builder for "upsert"-ing
	//  one ManifestTagReference node.
	ManifestTagReferenceUpsertOne struct {
		create *ManifestTagReferenceCreate
	}

	// ManifestTagReferenceUpsert is the "OnConflict" setter.
	ManifestTagReferenceUpsert struct {
		*sql.UpdateSet
	}
)

// SetTag sets the "tag" field.
func (u *ManifestTagReferenceUpsert) SetTag(v string) *ManifestTagReferenceUpsert {
	u.Set(manifesttagreference.FieldTag, v)
	return u
}

// UpdateTag sets the "tag" field to the value that was provided on create.
func (u *ManifestTagReferenceUpsert) UpdateTag() *ManifestTagReferenceUpsert {
	u.SetExcluded(manifesttagreference.FieldTag)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.ManifestTagReference.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ManifestTagReferenceUpsertOne) UpdateNewValues() *ManifestTagReferenceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ManifestTagReference.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ManifestTagReferenceUpsertOne) Ignore() *ManifestTagReferenceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ManifestTagReferenceUpsertOne) DoNothing() *ManifestTagReferenceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ManifestTagReferenceCreate.OnConflict
// documentation for more info.
func (u *ManifestTagReferenceUpsertOne) Update(set func(*ManifestTagReferenceUpsert)) *ManifestTagReferenceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ManifestTagReferenceUpsert{UpdateSet: update})
	}))
	return u
}

// SetTag sets the "tag" field.
func (u *ManifestTagReferenceUpsertOne) SetTag(v string) *ManifestTagReferenceUpsertOne {
	return u.Update(func(s *ManifestTagReferenceUpsert) {
		s.SetTag(v)
	})
}

// UpdateTag sets the "tag" field to the value that was provided on create.
func (u *ManifestTagReferenceUpsertOne) UpdateTag() *ManifestTagReferenceUpsertOne {
	return u.Update(func(s *ManifestTagReferenceUpsert) {
		s.UpdateTag()
	})
}

// Exec executes the query.
func (u *ManifestTagReferenceUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ManifestTagReferenceCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ManifestTagReferenceUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ManifestTagReferenceUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ManifestTagReferenceUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ManifestTagReferenceCreateBulk is the builder for creating many ManifestTagReference entities in bulk.
type ManifestTagReferenceCreateBulk struct {
	config
	err      error
	builders []*ManifestTagReferenceCreate
	conflict []sql.ConflictOption
}

// Save creates the ManifestTagReference entities in the database.
func (mtrcb *ManifestTagReferenceCreateBulk) Save(ctx context.Context) ([]*ManifestTagReference, error) {
	if mtrcb.err != nil {
		return nil, mtrcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mtrcb.builders))
	nodes := make([]*ManifestTagReference, len(mtrcb.builders))
	mutators := make([]Mutator, len(mtrcb.builders))
	for i := range mtrcb.builders {
		func(i int, root context.Context) {
			builder := mtrcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ManifestTagReferenceMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, mtrcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = mtrcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mtrcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, mtrcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mtrcb *ManifestTagReferenceCreateBulk) SaveX(ctx context.Context) []*ManifestTagReference {
	v, err := mtrcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mtrcb *ManifestTagReferenceCreateBulk) Exec(ctx context.Context) error {
	_, err := mtrcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mtrcb *ManifestTagReferenceCreateBulk) ExecX(ctx context.Context) {
	if err := mtrcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ManifestTagReference.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ManifestTagReferenceUpsert) {
//			SetTag(v+v).
//		}).
//		Exec(ctx)
func (mtrcb *ManifestTagReferenceCreateBulk) OnConflict(opts ...sql.ConflictOption) *ManifestTagReferenceUpsertBulk {
	mtrcb.conflict = opts
	return &ManifestTagReferenceUpsertBulk{
		create: mtrcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ManifestTagReference.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mtrcb *ManifestTagReferenceCreateBulk) OnConflictColumns(columns ...string) *ManifestTagReferenceUpsertBulk {
	mtrcb.conflict = append(mtrcb.conflict, sql.ConflictColumns(columns...))
	return &ManifestTagReferenceUpsertBulk{
		create: mtrcb,
	}
}

// ManifestTagReferenceUpsertBulk is the builder for "upsert"-ing
// a bulk of ManifestTagReference nodes.
type ManifestTagReferenceUpsertBulk struct {
	create *ManifestTagReferenceCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ManifestTagReference.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ManifestTagReferenceUpsertBulk) UpdateNewValues() *ManifestTagReferenceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ManifestTagReference.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ManifestTagReferenceUpsertBulk) Ignore() *ManifestTagReferenceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ManifestTagReferenceUpsertBulk) DoNothing() *ManifestTagReferenceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ManifestTagReferenceCreateBulk.OnConflict
// documentation for more info.
func (u *ManifestTagReferenceUpsertBulk) Update(set func(*ManifestTagReferenceUpsert)) *ManifestTagReferenceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ManifestTagReferenceUpsert{UpdateSet: update})
	}))
	return u
}

// SetTag sets the "tag" field.
func (u *ManifestTagReferenceUpsertBulk) SetTag(v string) *ManifestTagReferenceUpsertBulk {
	return u.Update(func(s *ManifestTagReferenceUpsert) {
		s.SetTag(v)
	})
}

// UpdateTag sets the "tag" field to the value that was provided on create.
func (u *ManifestTagReferenceUpsertBulk) UpdateTag() *ManifestTagReferenceUpsertBulk {
	return u.Update(func(s *ManifestTagReferenceUpsert) {
		s.UpdateTag()
	})
}

// Exec executes the query.
func (u *ManifestTagReferenceUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the ManifestTagReferenceCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ManifestTagReferenceCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ManifestTagReferenceUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
