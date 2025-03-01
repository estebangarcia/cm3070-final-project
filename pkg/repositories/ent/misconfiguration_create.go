// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifestmisconfiguration"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/misconfiguration"
)

// MisconfigurationCreate is the builder for creating a Misconfiguration entity.
type MisconfigurationCreate struct {
	config
	mutation *MisconfigurationMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetMisconfigurationID sets the "misconfiguration_id" field.
func (mc *MisconfigurationCreate) SetMisconfigurationID(s string) *MisconfigurationCreate {
	mc.mutation.SetMisconfigurationID(s)
	return mc
}

// SetMisconfigurationURLDetails sets the "misconfiguration_url_details" field.
func (mc *MisconfigurationCreate) SetMisconfigurationURLDetails(s string) *MisconfigurationCreate {
	mc.mutation.SetMisconfigurationURLDetails(s)
	return mc
}

// SetTitle sets the "title" field.
func (mc *MisconfigurationCreate) SetTitle(s string) *MisconfigurationCreate {
	mc.mutation.SetTitle(s)
	return mc
}

// SetSeverity sets the "severity" field.
func (mc *MisconfigurationCreate) SetSeverity(m misconfiguration.Severity) *MisconfigurationCreate {
	mc.mutation.SetSeverity(m)
	return mc
}

// AddManifestMisconfigurationIDs adds the "manifest_misconfigurations" edge to the ManifestMisconfiguration entity by IDs.
func (mc *MisconfigurationCreate) AddManifestMisconfigurationIDs(ids ...int) *MisconfigurationCreate {
	mc.mutation.AddManifestMisconfigurationIDs(ids...)
	return mc
}

// AddManifestMisconfigurations adds the "manifest_misconfigurations" edges to the ManifestMisconfiguration entity.
func (mc *MisconfigurationCreate) AddManifestMisconfigurations(m ...*ManifestMisconfiguration) *MisconfigurationCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mc.AddManifestMisconfigurationIDs(ids...)
}

// Mutation returns the MisconfigurationMutation object of the builder.
func (mc *MisconfigurationCreate) Mutation() *MisconfigurationMutation {
	return mc.mutation
}

// Save creates the Misconfiguration in the database.
func (mc *MisconfigurationCreate) Save(ctx context.Context) (*Misconfiguration, error) {
	return withHooks(ctx, mc.sqlSave, mc.mutation, mc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MisconfigurationCreate) SaveX(ctx context.Context) *Misconfiguration {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MisconfigurationCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MisconfigurationCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MisconfigurationCreate) check() error {
	if _, ok := mc.mutation.MisconfigurationID(); !ok {
		return &ValidationError{Name: "misconfiguration_id", err: errors.New(`ent: missing required field "Misconfiguration.misconfiguration_id"`)}
	}
	if _, ok := mc.mutation.MisconfigurationURLDetails(); !ok {
		return &ValidationError{Name: "misconfiguration_url_details", err: errors.New(`ent: missing required field "Misconfiguration.misconfiguration_url_details"`)}
	}
	if _, ok := mc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "Misconfiguration.title"`)}
	}
	if _, ok := mc.mutation.Severity(); !ok {
		return &ValidationError{Name: "severity", err: errors.New(`ent: missing required field "Misconfiguration.severity"`)}
	}
	if v, ok := mc.mutation.Severity(); ok {
		if err := misconfiguration.SeverityValidator(v); err != nil {
			return &ValidationError{Name: "severity", err: fmt.Errorf(`ent: validator failed for field "Misconfiguration.severity": %w`, err)}
		}
	}
	return nil
}

func (mc *MisconfigurationCreate) sqlSave(ctx context.Context) (*Misconfiguration, error) {
	if err := mc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	mc.mutation.id = &_node.ID
	mc.mutation.done = true
	return _node, nil
}

func (mc *MisconfigurationCreate) createSpec() (*Misconfiguration, *sqlgraph.CreateSpec) {
	var (
		_node = &Misconfiguration{config: mc.config}
		_spec = sqlgraph.NewCreateSpec(misconfiguration.Table, sqlgraph.NewFieldSpec(misconfiguration.FieldID, field.TypeInt))
	)
	_spec.OnConflict = mc.conflict
	if value, ok := mc.mutation.MisconfigurationID(); ok {
		_spec.SetField(misconfiguration.FieldMisconfigurationID, field.TypeString, value)
		_node.MisconfigurationID = value
	}
	if value, ok := mc.mutation.MisconfigurationURLDetails(); ok {
		_spec.SetField(misconfiguration.FieldMisconfigurationURLDetails, field.TypeString, value)
		_node.MisconfigurationURLDetails = value
	}
	if value, ok := mc.mutation.Title(); ok {
		_spec.SetField(misconfiguration.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := mc.mutation.Severity(); ok {
		_spec.SetField(misconfiguration.FieldSeverity, field.TypeEnum, value)
		_node.Severity = value
	}
	if nodes := mc.mutation.ManifestMisconfigurationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   misconfiguration.ManifestMisconfigurationsTable,
			Columns: []string{misconfiguration.ManifestMisconfigurationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifestmisconfiguration.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Misconfiguration.Create().
//		SetMisconfigurationID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MisconfigurationUpsert) {
//			SetMisconfigurationID(v+v).
//		}).
//		Exec(ctx)
func (mc *MisconfigurationCreate) OnConflict(opts ...sql.ConflictOption) *MisconfigurationUpsertOne {
	mc.conflict = opts
	return &MisconfigurationUpsertOne{
		create: mc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Misconfiguration.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mc *MisconfigurationCreate) OnConflictColumns(columns ...string) *MisconfigurationUpsertOne {
	mc.conflict = append(mc.conflict, sql.ConflictColumns(columns...))
	return &MisconfigurationUpsertOne{
		create: mc,
	}
}

type (
	// MisconfigurationUpsertOne is the builder for "upsert"-ing
	//  one Misconfiguration node.
	MisconfigurationUpsertOne struct {
		create *MisconfigurationCreate
	}

	// MisconfigurationUpsert is the "OnConflict" setter.
	MisconfigurationUpsert struct {
		*sql.UpdateSet
	}
)

// SetMisconfigurationID sets the "misconfiguration_id" field.
func (u *MisconfigurationUpsert) SetMisconfigurationID(v string) *MisconfigurationUpsert {
	u.Set(misconfiguration.FieldMisconfigurationID, v)
	return u
}

// UpdateMisconfigurationID sets the "misconfiguration_id" field to the value that was provided on create.
func (u *MisconfigurationUpsert) UpdateMisconfigurationID() *MisconfigurationUpsert {
	u.SetExcluded(misconfiguration.FieldMisconfigurationID)
	return u
}

// SetMisconfigurationURLDetails sets the "misconfiguration_url_details" field.
func (u *MisconfigurationUpsert) SetMisconfigurationURLDetails(v string) *MisconfigurationUpsert {
	u.Set(misconfiguration.FieldMisconfigurationURLDetails, v)
	return u
}

// UpdateMisconfigurationURLDetails sets the "misconfiguration_url_details" field to the value that was provided on create.
func (u *MisconfigurationUpsert) UpdateMisconfigurationURLDetails() *MisconfigurationUpsert {
	u.SetExcluded(misconfiguration.FieldMisconfigurationURLDetails)
	return u
}

// SetTitle sets the "title" field.
func (u *MisconfigurationUpsert) SetTitle(v string) *MisconfigurationUpsert {
	u.Set(misconfiguration.FieldTitle, v)
	return u
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *MisconfigurationUpsert) UpdateTitle() *MisconfigurationUpsert {
	u.SetExcluded(misconfiguration.FieldTitle)
	return u
}

// SetSeverity sets the "severity" field.
func (u *MisconfigurationUpsert) SetSeverity(v misconfiguration.Severity) *MisconfigurationUpsert {
	u.Set(misconfiguration.FieldSeverity, v)
	return u
}

// UpdateSeverity sets the "severity" field to the value that was provided on create.
func (u *MisconfigurationUpsert) UpdateSeverity() *MisconfigurationUpsert {
	u.SetExcluded(misconfiguration.FieldSeverity)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Misconfiguration.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *MisconfigurationUpsertOne) UpdateNewValues() *MisconfigurationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Misconfiguration.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *MisconfigurationUpsertOne) Ignore() *MisconfigurationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MisconfigurationUpsertOne) DoNothing() *MisconfigurationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MisconfigurationCreate.OnConflict
// documentation for more info.
func (u *MisconfigurationUpsertOne) Update(set func(*MisconfigurationUpsert)) *MisconfigurationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MisconfigurationUpsert{UpdateSet: update})
	}))
	return u
}

// SetMisconfigurationID sets the "misconfiguration_id" field.
func (u *MisconfigurationUpsertOne) SetMisconfigurationID(v string) *MisconfigurationUpsertOne {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.SetMisconfigurationID(v)
	})
}

// UpdateMisconfigurationID sets the "misconfiguration_id" field to the value that was provided on create.
func (u *MisconfigurationUpsertOne) UpdateMisconfigurationID() *MisconfigurationUpsertOne {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.UpdateMisconfigurationID()
	})
}

// SetMisconfigurationURLDetails sets the "misconfiguration_url_details" field.
func (u *MisconfigurationUpsertOne) SetMisconfigurationURLDetails(v string) *MisconfigurationUpsertOne {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.SetMisconfigurationURLDetails(v)
	})
}

// UpdateMisconfigurationURLDetails sets the "misconfiguration_url_details" field to the value that was provided on create.
func (u *MisconfigurationUpsertOne) UpdateMisconfigurationURLDetails() *MisconfigurationUpsertOne {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.UpdateMisconfigurationURLDetails()
	})
}

// SetTitle sets the "title" field.
func (u *MisconfigurationUpsertOne) SetTitle(v string) *MisconfigurationUpsertOne {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *MisconfigurationUpsertOne) UpdateTitle() *MisconfigurationUpsertOne {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.UpdateTitle()
	})
}

// SetSeverity sets the "severity" field.
func (u *MisconfigurationUpsertOne) SetSeverity(v misconfiguration.Severity) *MisconfigurationUpsertOne {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.SetSeverity(v)
	})
}

// UpdateSeverity sets the "severity" field to the value that was provided on create.
func (u *MisconfigurationUpsertOne) UpdateSeverity() *MisconfigurationUpsertOne {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.UpdateSeverity()
	})
}

// Exec executes the query.
func (u *MisconfigurationUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for MisconfigurationCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MisconfigurationUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *MisconfigurationUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *MisconfigurationUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// MisconfigurationCreateBulk is the builder for creating many Misconfiguration entities in bulk.
type MisconfigurationCreateBulk struct {
	config
	err      error
	builders []*MisconfigurationCreate
	conflict []sql.ConflictOption
}

// Save creates the Misconfiguration entities in the database.
func (mcb *MisconfigurationCreateBulk) Save(ctx context.Context) ([]*Misconfiguration, error) {
	if mcb.err != nil {
		return nil, mcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Misconfiguration, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MisconfigurationMutation)
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
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = mcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MisconfigurationCreateBulk) SaveX(ctx context.Context) []*Misconfiguration {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MisconfigurationCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MisconfigurationCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Misconfiguration.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MisconfigurationUpsert) {
//			SetMisconfigurationID(v+v).
//		}).
//		Exec(ctx)
func (mcb *MisconfigurationCreateBulk) OnConflict(opts ...sql.ConflictOption) *MisconfigurationUpsertBulk {
	mcb.conflict = opts
	return &MisconfigurationUpsertBulk{
		create: mcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Misconfiguration.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mcb *MisconfigurationCreateBulk) OnConflictColumns(columns ...string) *MisconfigurationUpsertBulk {
	mcb.conflict = append(mcb.conflict, sql.ConflictColumns(columns...))
	return &MisconfigurationUpsertBulk{
		create: mcb,
	}
}

// MisconfigurationUpsertBulk is the builder for "upsert"-ing
// a bulk of Misconfiguration nodes.
type MisconfigurationUpsertBulk struct {
	create *MisconfigurationCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Misconfiguration.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *MisconfigurationUpsertBulk) UpdateNewValues() *MisconfigurationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Misconfiguration.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *MisconfigurationUpsertBulk) Ignore() *MisconfigurationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MisconfigurationUpsertBulk) DoNothing() *MisconfigurationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MisconfigurationCreateBulk.OnConflict
// documentation for more info.
func (u *MisconfigurationUpsertBulk) Update(set func(*MisconfigurationUpsert)) *MisconfigurationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MisconfigurationUpsert{UpdateSet: update})
	}))
	return u
}

// SetMisconfigurationID sets the "misconfiguration_id" field.
func (u *MisconfigurationUpsertBulk) SetMisconfigurationID(v string) *MisconfigurationUpsertBulk {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.SetMisconfigurationID(v)
	})
}

// UpdateMisconfigurationID sets the "misconfiguration_id" field to the value that was provided on create.
func (u *MisconfigurationUpsertBulk) UpdateMisconfigurationID() *MisconfigurationUpsertBulk {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.UpdateMisconfigurationID()
	})
}

// SetMisconfigurationURLDetails sets the "misconfiguration_url_details" field.
func (u *MisconfigurationUpsertBulk) SetMisconfigurationURLDetails(v string) *MisconfigurationUpsertBulk {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.SetMisconfigurationURLDetails(v)
	})
}

// UpdateMisconfigurationURLDetails sets the "misconfiguration_url_details" field to the value that was provided on create.
func (u *MisconfigurationUpsertBulk) UpdateMisconfigurationURLDetails() *MisconfigurationUpsertBulk {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.UpdateMisconfigurationURLDetails()
	})
}

// SetTitle sets the "title" field.
func (u *MisconfigurationUpsertBulk) SetTitle(v string) *MisconfigurationUpsertBulk {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *MisconfigurationUpsertBulk) UpdateTitle() *MisconfigurationUpsertBulk {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.UpdateTitle()
	})
}

// SetSeverity sets the "severity" field.
func (u *MisconfigurationUpsertBulk) SetSeverity(v misconfiguration.Severity) *MisconfigurationUpsertBulk {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.SetSeverity(v)
	})
}

// UpdateSeverity sets the "severity" field to the value that was provided on create.
func (u *MisconfigurationUpsertBulk) UpdateSeverity() *MisconfigurationUpsertBulk {
	return u.Update(func(s *MisconfigurationUpsert) {
		s.UpdateSeverity()
	})
}

// Exec executes the query.
func (u *MisconfigurationUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the MisconfigurationCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for MisconfigurationCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MisconfigurationUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
