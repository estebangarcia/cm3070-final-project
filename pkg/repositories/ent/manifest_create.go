// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifesttagreference"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
)

// ManifestCreate is the builder for creating a Manifest entity.
type ManifestCreate struct {
	config
	mutation *ManifestMutation
	hooks    []Hook
}

// SetMediaType sets the "media_type" field.
func (mc *ManifestCreate) SetMediaType(s string) *ManifestCreate {
	mc.mutation.SetMediaType(s)
	return mc
}

// SetS3Path sets the "s3_path" field.
func (mc *ManifestCreate) SetS3Path(s string) *ManifestCreate {
	mc.mutation.SetS3Path(s)
	return mc
}

// SetDigest sets the "digest" field.
func (mc *ManifestCreate) SetDigest(s string) *ManifestCreate {
	mc.mutation.SetDigest(s)
	return mc
}

// AddTagIDs adds the "tags" edge to the ManifestTagReference entity by IDs.
func (mc *ManifestCreate) AddTagIDs(ids ...int) *ManifestCreate {
	mc.mutation.AddTagIDs(ids...)
	return mc
}

// AddTags adds the "tags" edges to the ManifestTagReference entity.
func (mc *ManifestCreate) AddTags(m ...*ManifestTagReference) *ManifestCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mc.AddTagIDs(ids...)
}

// SetRepositoryID sets the "repository" edge to the Repository entity by ID.
func (mc *ManifestCreate) SetRepositoryID(id int) *ManifestCreate {
	mc.mutation.SetRepositoryID(id)
	return mc
}

// SetNillableRepositoryID sets the "repository" edge to the Repository entity by ID if the given value is not nil.
func (mc *ManifestCreate) SetNillableRepositoryID(id *int) *ManifestCreate {
	if id != nil {
		mc = mc.SetRepositoryID(*id)
	}
	return mc
}

// SetRepository sets the "repository" edge to the Repository entity.
func (mc *ManifestCreate) SetRepository(r *Repository) *ManifestCreate {
	return mc.SetRepositoryID(r.ID)
}

// AddSubjectIDs adds the "subject" edge to the Manifest entity by IDs.
func (mc *ManifestCreate) AddSubjectIDs(ids ...int) *ManifestCreate {
	mc.mutation.AddSubjectIDs(ids...)
	return mc
}

// AddSubject adds the "subject" edges to the Manifest entity.
func (mc *ManifestCreate) AddSubject(m ...*Manifest) *ManifestCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mc.AddSubjectIDs(ids...)
}

// AddRefererIDs adds the "referer" edge to the Manifest entity by IDs.
func (mc *ManifestCreate) AddRefererIDs(ids ...int) *ManifestCreate {
	mc.mutation.AddRefererIDs(ids...)
	return mc
}

// AddReferer adds the "referer" edges to the Manifest entity.
func (mc *ManifestCreate) AddReferer(m ...*Manifest) *ManifestCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mc.AddRefererIDs(ids...)
}

// Mutation returns the ManifestMutation object of the builder.
func (mc *ManifestCreate) Mutation() *ManifestMutation {
	return mc.mutation
}

// Save creates the Manifest in the database.
func (mc *ManifestCreate) Save(ctx context.Context) (*Manifest, error) {
	return withHooks(ctx, mc.sqlSave, mc.mutation, mc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mc *ManifestCreate) SaveX(ctx context.Context) *Manifest {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *ManifestCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *ManifestCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *ManifestCreate) check() error {
	if _, ok := mc.mutation.MediaType(); !ok {
		return &ValidationError{Name: "media_type", err: errors.New(`ent: missing required field "Manifest.media_type"`)}
	}
	if _, ok := mc.mutation.S3Path(); !ok {
		return &ValidationError{Name: "s3_path", err: errors.New(`ent: missing required field "Manifest.s3_path"`)}
	}
	if _, ok := mc.mutation.Digest(); !ok {
		return &ValidationError{Name: "digest", err: errors.New(`ent: missing required field "Manifest.digest"`)}
	}
	return nil
}

func (mc *ManifestCreate) sqlSave(ctx context.Context) (*Manifest, error) {
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

func (mc *ManifestCreate) createSpec() (*Manifest, *sqlgraph.CreateSpec) {
	var (
		_node = &Manifest{config: mc.config}
		_spec = sqlgraph.NewCreateSpec(manifest.Table, sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt))
	)
	if value, ok := mc.mutation.MediaType(); ok {
		_spec.SetField(manifest.FieldMediaType, field.TypeString, value)
		_node.MediaType = value
	}
	if value, ok := mc.mutation.S3Path(); ok {
		_spec.SetField(manifest.FieldS3Path, field.TypeString, value)
		_node.S3Path = value
	}
	if value, ok := mc.mutation.Digest(); ok {
		_spec.SetField(manifest.FieldDigest, field.TypeString, value)
		_node.Digest = value
	}
	if nodes := mc.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   manifest.TagsTable,
			Columns: []string{manifest.TagsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifesttagreference.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.RepositoryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   manifest.RepositoryTable,
			Columns: []string{manifest.RepositoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(repository.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.repository_manifests = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.SubjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   manifest.SubjectTable,
			Columns: manifest.SubjectPrimaryKey,
			Bidi:    true,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.RefererIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   manifest.RefererTable,
			Columns: manifest.RefererPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ManifestCreateBulk is the builder for creating many Manifest entities in bulk.
type ManifestCreateBulk struct {
	config
	err      error
	builders []*ManifestCreate
}

// Save creates the Manifest entities in the database.
func (mcb *ManifestCreateBulk) Save(ctx context.Context) ([]*Manifest, error) {
	if mcb.err != nil {
		return nil, mcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Manifest, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ManifestMutation)
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
func (mcb *ManifestCreateBulk) SaveX(ctx context.Context) []*Manifest {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *ManifestCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *ManifestCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}
