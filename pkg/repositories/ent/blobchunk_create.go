// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/blobchunk"
)

// BlobChunkCreate is the builder for creating a BlobChunk entity.
type BlobChunkCreate struct {
	config
	mutation *BlobChunkMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetUploadID sets the "upload_id" field.
func (bcc *BlobChunkCreate) SetUploadID(s string) *BlobChunkCreate {
	bcc.mutation.SetUploadID(s)
	return bcc
}

// SetSessionID sets the "session_id" field.
func (bcc *BlobChunkCreate) SetSessionID(s string) *BlobChunkCreate {
	bcc.mutation.SetSessionID(s)
	return bcc
}

// SetRangeFrom sets the "range_from" field.
func (bcc *BlobChunkCreate) SetRangeFrom(u uint64) *BlobChunkCreate {
	bcc.mutation.SetRangeFrom(u)
	return bcc
}

// SetRangeTo sets the "range_to" field.
func (bcc *BlobChunkCreate) SetRangeTo(u uint64) *BlobChunkCreate {
	bcc.mutation.SetRangeTo(u)
	return bcc
}

// SetPartNumber sets the "part_number" field.
func (bcc *BlobChunkCreate) SetPartNumber(u uint64) *BlobChunkCreate {
	bcc.mutation.SetPartNumber(u)
	return bcc
}

// Mutation returns the BlobChunkMutation object of the builder.
func (bcc *BlobChunkCreate) Mutation() *BlobChunkMutation {
	return bcc.mutation
}

// Save creates the BlobChunk in the database.
func (bcc *BlobChunkCreate) Save(ctx context.Context) (*BlobChunk, error) {
	return withHooks(ctx, bcc.sqlSave, bcc.mutation, bcc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (bcc *BlobChunkCreate) SaveX(ctx context.Context) *BlobChunk {
	v, err := bcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bcc *BlobChunkCreate) Exec(ctx context.Context) error {
	_, err := bcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bcc *BlobChunkCreate) ExecX(ctx context.Context) {
	if err := bcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bcc *BlobChunkCreate) check() error {
	if _, ok := bcc.mutation.UploadID(); !ok {
		return &ValidationError{Name: "upload_id", err: errors.New(`ent: missing required field "BlobChunk.upload_id"`)}
	}
	if _, ok := bcc.mutation.SessionID(); !ok {
		return &ValidationError{Name: "session_id", err: errors.New(`ent: missing required field "BlobChunk.session_id"`)}
	}
	if _, ok := bcc.mutation.RangeFrom(); !ok {
		return &ValidationError{Name: "range_from", err: errors.New(`ent: missing required field "BlobChunk.range_from"`)}
	}
	if _, ok := bcc.mutation.RangeTo(); !ok {
		return &ValidationError{Name: "range_to", err: errors.New(`ent: missing required field "BlobChunk.range_to"`)}
	}
	if _, ok := bcc.mutation.PartNumber(); !ok {
		return &ValidationError{Name: "part_number", err: errors.New(`ent: missing required field "BlobChunk.part_number"`)}
	}
	return nil
}

func (bcc *BlobChunkCreate) sqlSave(ctx context.Context) (*BlobChunk, error) {
	if err := bcc.check(); err != nil {
		return nil, err
	}
	_node, _spec := bcc.createSpec()
	if err := sqlgraph.CreateNode(ctx, bcc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	bcc.mutation.id = &_node.ID
	bcc.mutation.done = true
	return _node, nil
}

func (bcc *BlobChunkCreate) createSpec() (*BlobChunk, *sqlgraph.CreateSpec) {
	var (
		_node = &BlobChunk{config: bcc.config}
		_spec = sqlgraph.NewCreateSpec(blobchunk.Table, sqlgraph.NewFieldSpec(blobchunk.FieldID, field.TypeInt))
	)
	_spec.OnConflict = bcc.conflict
	if value, ok := bcc.mutation.UploadID(); ok {
		_spec.SetField(blobchunk.FieldUploadID, field.TypeString, value)
		_node.UploadID = value
	}
	if value, ok := bcc.mutation.SessionID(); ok {
		_spec.SetField(blobchunk.FieldSessionID, field.TypeString, value)
		_node.SessionID = value
	}
	if value, ok := bcc.mutation.RangeFrom(); ok {
		_spec.SetField(blobchunk.FieldRangeFrom, field.TypeUint64, value)
		_node.RangeFrom = value
	}
	if value, ok := bcc.mutation.RangeTo(); ok {
		_spec.SetField(blobchunk.FieldRangeTo, field.TypeUint64, value)
		_node.RangeTo = value
	}
	if value, ok := bcc.mutation.PartNumber(); ok {
		_spec.SetField(blobchunk.FieldPartNumber, field.TypeUint64, value)
		_node.PartNumber = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.BlobChunk.Create().
//		SetUploadID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.BlobChunkUpsert) {
//			SetUploadID(v+v).
//		}).
//		Exec(ctx)
func (bcc *BlobChunkCreate) OnConflict(opts ...sql.ConflictOption) *BlobChunkUpsertOne {
	bcc.conflict = opts
	return &BlobChunkUpsertOne{
		create: bcc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.BlobChunk.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (bcc *BlobChunkCreate) OnConflictColumns(columns ...string) *BlobChunkUpsertOne {
	bcc.conflict = append(bcc.conflict, sql.ConflictColumns(columns...))
	return &BlobChunkUpsertOne{
		create: bcc,
	}
}

type (
	// BlobChunkUpsertOne is the builder for "upsert"-ing
	//  one BlobChunk node.
	BlobChunkUpsertOne struct {
		create *BlobChunkCreate
	}

	// BlobChunkUpsert is the "OnConflict" setter.
	BlobChunkUpsert struct {
		*sql.UpdateSet
	}
)

// SetUploadID sets the "upload_id" field.
func (u *BlobChunkUpsert) SetUploadID(v string) *BlobChunkUpsert {
	u.Set(blobchunk.FieldUploadID, v)
	return u
}

// UpdateUploadID sets the "upload_id" field to the value that was provided on create.
func (u *BlobChunkUpsert) UpdateUploadID() *BlobChunkUpsert {
	u.SetExcluded(blobchunk.FieldUploadID)
	return u
}

// SetSessionID sets the "session_id" field.
func (u *BlobChunkUpsert) SetSessionID(v string) *BlobChunkUpsert {
	u.Set(blobchunk.FieldSessionID, v)
	return u
}

// UpdateSessionID sets the "session_id" field to the value that was provided on create.
func (u *BlobChunkUpsert) UpdateSessionID() *BlobChunkUpsert {
	u.SetExcluded(blobchunk.FieldSessionID)
	return u
}

// SetRangeFrom sets the "range_from" field.
func (u *BlobChunkUpsert) SetRangeFrom(v uint64) *BlobChunkUpsert {
	u.Set(blobchunk.FieldRangeFrom, v)
	return u
}

// UpdateRangeFrom sets the "range_from" field to the value that was provided on create.
func (u *BlobChunkUpsert) UpdateRangeFrom() *BlobChunkUpsert {
	u.SetExcluded(blobchunk.FieldRangeFrom)
	return u
}

// AddRangeFrom adds v to the "range_from" field.
func (u *BlobChunkUpsert) AddRangeFrom(v uint64) *BlobChunkUpsert {
	u.Add(blobchunk.FieldRangeFrom, v)
	return u
}

// SetRangeTo sets the "range_to" field.
func (u *BlobChunkUpsert) SetRangeTo(v uint64) *BlobChunkUpsert {
	u.Set(blobchunk.FieldRangeTo, v)
	return u
}

// UpdateRangeTo sets the "range_to" field to the value that was provided on create.
func (u *BlobChunkUpsert) UpdateRangeTo() *BlobChunkUpsert {
	u.SetExcluded(blobchunk.FieldRangeTo)
	return u
}

// AddRangeTo adds v to the "range_to" field.
func (u *BlobChunkUpsert) AddRangeTo(v uint64) *BlobChunkUpsert {
	u.Add(blobchunk.FieldRangeTo, v)
	return u
}

// SetPartNumber sets the "part_number" field.
func (u *BlobChunkUpsert) SetPartNumber(v uint64) *BlobChunkUpsert {
	u.Set(blobchunk.FieldPartNumber, v)
	return u
}

// UpdatePartNumber sets the "part_number" field to the value that was provided on create.
func (u *BlobChunkUpsert) UpdatePartNumber() *BlobChunkUpsert {
	u.SetExcluded(blobchunk.FieldPartNumber)
	return u
}

// AddPartNumber adds v to the "part_number" field.
func (u *BlobChunkUpsert) AddPartNumber(v uint64) *BlobChunkUpsert {
	u.Add(blobchunk.FieldPartNumber, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.BlobChunk.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *BlobChunkUpsertOne) UpdateNewValues() *BlobChunkUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.BlobChunk.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *BlobChunkUpsertOne) Ignore() *BlobChunkUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *BlobChunkUpsertOne) DoNothing() *BlobChunkUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the BlobChunkCreate.OnConflict
// documentation for more info.
func (u *BlobChunkUpsertOne) Update(set func(*BlobChunkUpsert)) *BlobChunkUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&BlobChunkUpsert{UpdateSet: update})
	}))
	return u
}

// SetUploadID sets the "upload_id" field.
func (u *BlobChunkUpsertOne) SetUploadID(v string) *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetUploadID(v)
	})
}

// UpdateUploadID sets the "upload_id" field to the value that was provided on create.
func (u *BlobChunkUpsertOne) UpdateUploadID() *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdateUploadID()
	})
}

// SetSessionID sets the "session_id" field.
func (u *BlobChunkUpsertOne) SetSessionID(v string) *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetSessionID(v)
	})
}

// UpdateSessionID sets the "session_id" field to the value that was provided on create.
func (u *BlobChunkUpsertOne) UpdateSessionID() *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdateSessionID()
	})
}

// SetRangeFrom sets the "range_from" field.
func (u *BlobChunkUpsertOne) SetRangeFrom(v uint64) *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetRangeFrom(v)
	})
}

// AddRangeFrom adds v to the "range_from" field.
func (u *BlobChunkUpsertOne) AddRangeFrom(v uint64) *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.AddRangeFrom(v)
	})
}

// UpdateRangeFrom sets the "range_from" field to the value that was provided on create.
func (u *BlobChunkUpsertOne) UpdateRangeFrom() *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdateRangeFrom()
	})
}

// SetRangeTo sets the "range_to" field.
func (u *BlobChunkUpsertOne) SetRangeTo(v uint64) *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetRangeTo(v)
	})
}

// AddRangeTo adds v to the "range_to" field.
func (u *BlobChunkUpsertOne) AddRangeTo(v uint64) *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.AddRangeTo(v)
	})
}

// UpdateRangeTo sets the "range_to" field to the value that was provided on create.
func (u *BlobChunkUpsertOne) UpdateRangeTo() *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdateRangeTo()
	})
}

// SetPartNumber sets the "part_number" field.
func (u *BlobChunkUpsertOne) SetPartNumber(v uint64) *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetPartNumber(v)
	})
}

// AddPartNumber adds v to the "part_number" field.
func (u *BlobChunkUpsertOne) AddPartNumber(v uint64) *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.AddPartNumber(v)
	})
}

// UpdatePartNumber sets the "part_number" field to the value that was provided on create.
func (u *BlobChunkUpsertOne) UpdatePartNumber() *BlobChunkUpsertOne {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdatePartNumber()
	})
}

// Exec executes the query.
func (u *BlobChunkUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for BlobChunkCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *BlobChunkUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *BlobChunkUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *BlobChunkUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// BlobChunkCreateBulk is the builder for creating many BlobChunk entities in bulk.
type BlobChunkCreateBulk struct {
	config
	err      error
	builders []*BlobChunkCreate
	conflict []sql.ConflictOption
}

// Save creates the BlobChunk entities in the database.
func (bccb *BlobChunkCreateBulk) Save(ctx context.Context) ([]*BlobChunk, error) {
	if bccb.err != nil {
		return nil, bccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(bccb.builders))
	nodes := make([]*BlobChunk, len(bccb.builders))
	mutators := make([]Mutator, len(bccb.builders))
	for i := range bccb.builders {
		func(i int, root context.Context) {
			builder := bccb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*BlobChunkMutation)
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
					_, err = mutators[i+1].Mutate(root, bccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = bccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, bccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, bccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (bccb *BlobChunkCreateBulk) SaveX(ctx context.Context) []*BlobChunk {
	v, err := bccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bccb *BlobChunkCreateBulk) Exec(ctx context.Context) error {
	_, err := bccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bccb *BlobChunkCreateBulk) ExecX(ctx context.Context) {
	if err := bccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.BlobChunk.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.BlobChunkUpsert) {
//			SetUploadID(v+v).
//		}).
//		Exec(ctx)
func (bccb *BlobChunkCreateBulk) OnConflict(opts ...sql.ConflictOption) *BlobChunkUpsertBulk {
	bccb.conflict = opts
	return &BlobChunkUpsertBulk{
		create: bccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.BlobChunk.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (bccb *BlobChunkCreateBulk) OnConflictColumns(columns ...string) *BlobChunkUpsertBulk {
	bccb.conflict = append(bccb.conflict, sql.ConflictColumns(columns...))
	return &BlobChunkUpsertBulk{
		create: bccb,
	}
}

// BlobChunkUpsertBulk is the builder for "upsert"-ing
// a bulk of BlobChunk nodes.
type BlobChunkUpsertBulk struct {
	create *BlobChunkCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.BlobChunk.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *BlobChunkUpsertBulk) UpdateNewValues() *BlobChunkUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.BlobChunk.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *BlobChunkUpsertBulk) Ignore() *BlobChunkUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *BlobChunkUpsertBulk) DoNothing() *BlobChunkUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the BlobChunkCreateBulk.OnConflict
// documentation for more info.
func (u *BlobChunkUpsertBulk) Update(set func(*BlobChunkUpsert)) *BlobChunkUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&BlobChunkUpsert{UpdateSet: update})
	}))
	return u
}

// SetUploadID sets the "upload_id" field.
func (u *BlobChunkUpsertBulk) SetUploadID(v string) *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetUploadID(v)
	})
}

// UpdateUploadID sets the "upload_id" field to the value that was provided on create.
func (u *BlobChunkUpsertBulk) UpdateUploadID() *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdateUploadID()
	})
}

// SetSessionID sets the "session_id" field.
func (u *BlobChunkUpsertBulk) SetSessionID(v string) *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetSessionID(v)
	})
}

// UpdateSessionID sets the "session_id" field to the value that was provided on create.
func (u *BlobChunkUpsertBulk) UpdateSessionID() *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdateSessionID()
	})
}

// SetRangeFrom sets the "range_from" field.
func (u *BlobChunkUpsertBulk) SetRangeFrom(v uint64) *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetRangeFrom(v)
	})
}

// AddRangeFrom adds v to the "range_from" field.
func (u *BlobChunkUpsertBulk) AddRangeFrom(v uint64) *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.AddRangeFrom(v)
	})
}

// UpdateRangeFrom sets the "range_from" field to the value that was provided on create.
func (u *BlobChunkUpsertBulk) UpdateRangeFrom() *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdateRangeFrom()
	})
}

// SetRangeTo sets the "range_to" field.
func (u *BlobChunkUpsertBulk) SetRangeTo(v uint64) *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetRangeTo(v)
	})
}

// AddRangeTo adds v to the "range_to" field.
func (u *BlobChunkUpsertBulk) AddRangeTo(v uint64) *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.AddRangeTo(v)
	})
}

// UpdateRangeTo sets the "range_to" field to the value that was provided on create.
func (u *BlobChunkUpsertBulk) UpdateRangeTo() *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdateRangeTo()
	})
}

// SetPartNumber sets the "part_number" field.
func (u *BlobChunkUpsertBulk) SetPartNumber(v uint64) *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.SetPartNumber(v)
	})
}

// AddPartNumber adds v to the "part_number" field.
func (u *BlobChunkUpsertBulk) AddPartNumber(v uint64) *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.AddPartNumber(v)
	})
}

// UpdatePartNumber sets the "part_number" field to the value that was provided on create.
func (u *BlobChunkUpsertBulk) UpdatePartNumber() *BlobChunkUpsertBulk {
	return u.Update(func(s *BlobChunkUpsert) {
		s.UpdatePartNumber()
	})
}

// Exec executes the query.
func (u *BlobChunkUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the BlobChunkCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for BlobChunkCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *BlobChunkUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
