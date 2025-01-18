// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifesttagreference"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
)

// ManifestQuery is the builder for querying Manifest entities.
type ManifestQuery struct {
	config
	ctx            *QueryContext
	order          []manifest.OrderOption
	inters         []Interceptor
	predicates     []predicate.Manifest
	withTags       *ManifestTagReferenceQuery
	withRepository *RepositoryQuery
	withSubject    *ManifestQuery
	withReferer    *ManifestQuery
	withFKs        bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ManifestQuery builder.
func (mq *ManifestQuery) Where(ps ...predicate.Manifest) *ManifestQuery {
	mq.predicates = append(mq.predicates, ps...)
	return mq
}

// Limit the number of records to be returned by this query.
func (mq *ManifestQuery) Limit(limit int) *ManifestQuery {
	mq.ctx.Limit = &limit
	return mq
}

// Offset to start from.
func (mq *ManifestQuery) Offset(offset int) *ManifestQuery {
	mq.ctx.Offset = &offset
	return mq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mq *ManifestQuery) Unique(unique bool) *ManifestQuery {
	mq.ctx.Unique = &unique
	return mq
}

// Order specifies how the records should be ordered.
func (mq *ManifestQuery) Order(o ...manifest.OrderOption) *ManifestQuery {
	mq.order = append(mq.order, o...)
	return mq
}

// QueryTags chains the current query on the "tags" edge.
func (mq *ManifestQuery) QueryTags() *ManifestTagReferenceQuery {
	query := (&ManifestTagReferenceClient{config: mq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(manifest.Table, manifest.FieldID, selector),
			sqlgraph.To(manifesttagreference.Table, manifesttagreference.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, manifest.TagsTable, manifest.TagsColumn),
		)
		fromU = sqlgraph.SetNeighbors(mq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRepository chains the current query on the "repository" edge.
func (mq *ManifestQuery) QueryRepository() *RepositoryQuery {
	query := (&RepositoryClient{config: mq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(manifest.Table, manifest.FieldID, selector),
			sqlgraph.To(repository.Table, repository.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, manifest.RepositoryTable, manifest.RepositoryColumn),
		)
		fromU = sqlgraph.SetNeighbors(mq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySubject chains the current query on the "subject" edge.
func (mq *ManifestQuery) QuerySubject() *ManifestQuery {
	query := (&ManifestClient{config: mq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(manifest.Table, manifest.FieldID, selector),
			sqlgraph.To(manifest.Table, manifest.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, manifest.SubjectTable, manifest.SubjectPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(mq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryReferer chains the current query on the "referer" edge.
func (mq *ManifestQuery) QueryReferer() *ManifestQuery {
	query := (&ManifestClient{config: mq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(manifest.Table, manifest.FieldID, selector),
			sqlgraph.To(manifest.Table, manifest.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, manifest.RefererTable, manifest.RefererPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(mq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Manifest entity from the query.
// Returns a *NotFoundError when no Manifest was found.
func (mq *ManifestQuery) First(ctx context.Context) (*Manifest, error) {
	nodes, err := mq.Limit(1).All(setContextOp(ctx, mq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{manifest.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mq *ManifestQuery) FirstX(ctx context.Context) *Manifest {
	node, err := mq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Manifest ID from the query.
// Returns a *NotFoundError when no Manifest ID was found.
func (mq *ManifestQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = mq.Limit(1).IDs(setContextOp(ctx, mq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{manifest.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mq *ManifestQuery) FirstIDX(ctx context.Context) int {
	id, err := mq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Manifest entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Manifest entity is found.
// Returns a *NotFoundError when no Manifest entities are found.
func (mq *ManifestQuery) Only(ctx context.Context) (*Manifest, error) {
	nodes, err := mq.Limit(2).All(setContextOp(ctx, mq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{manifest.Label}
	default:
		return nil, &NotSingularError{manifest.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mq *ManifestQuery) OnlyX(ctx context.Context) *Manifest {
	node, err := mq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Manifest ID in the query.
// Returns a *NotSingularError when more than one Manifest ID is found.
// Returns a *NotFoundError when no entities are found.
func (mq *ManifestQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = mq.Limit(2).IDs(setContextOp(ctx, mq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{manifest.Label}
	default:
		err = &NotSingularError{manifest.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mq *ManifestQuery) OnlyIDX(ctx context.Context) int {
	id, err := mq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Manifests.
func (mq *ManifestQuery) All(ctx context.Context) ([]*Manifest, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryAll)
	if err := mq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Manifest, *ManifestQuery]()
	return withInterceptors[[]*Manifest](ctx, mq, qr, mq.inters)
}

// AllX is like All, but panics if an error occurs.
func (mq *ManifestQuery) AllX(ctx context.Context) []*Manifest {
	nodes, err := mq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Manifest IDs.
func (mq *ManifestQuery) IDs(ctx context.Context) (ids []int, err error) {
	if mq.ctx.Unique == nil && mq.path != nil {
		mq.Unique(true)
	}
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryIDs)
	if err = mq.Select(manifest.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mq *ManifestQuery) IDsX(ctx context.Context) []int {
	ids, err := mq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mq *ManifestQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryCount)
	if err := mq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, mq, querierCount[*ManifestQuery](), mq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (mq *ManifestQuery) CountX(ctx context.Context) int {
	count, err := mq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mq *ManifestQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, mq.ctx, ent.OpQueryExist)
	switch _, err := mq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (mq *ManifestQuery) ExistX(ctx context.Context) bool {
	exist, err := mq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ManifestQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mq *ManifestQuery) Clone() *ManifestQuery {
	if mq == nil {
		return nil
	}
	return &ManifestQuery{
		config:         mq.config,
		ctx:            mq.ctx.Clone(),
		order:          append([]manifest.OrderOption{}, mq.order...),
		inters:         append([]Interceptor{}, mq.inters...),
		predicates:     append([]predicate.Manifest{}, mq.predicates...),
		withTags:       mq.withTags.Clone(),
		withRepository: mq.withRepository.Clone(),
		withSubject:    mq.withSubject.Clone(),
		withReferer:    mq.withReferer.Clone(),
		// clone intermediate query.
		sql:  mq.sql.Clone(),
		path: mq.path,
	}
}

// WithTags tells the query-builder to eager-load the nodes that are connected to
// the "tags" edge. The optional arguments are used to configure the query builder of the edge.
func (mq *ManifestQuery) WithTags(opts ...func(*ManifestTagReferenceQuery)) *ManifestQuery {
	query := (&ManifestTagReferenceClient{config: mq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mq.withTags = query
	return mq
}

// WithRepository tells the query-builder to eager-load the nodes that are connected to
// the "repository" edge. The optional arguments are used to configure the query builder of the edge.
func (mq *ManifestQuery) WithRepository(opts ...func(*RepositoryQuery)) *ManifestQuery {
	query := (&RepositoryClient{config: mq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mq.withRepository = query
	return mq
}

// WithSubject tells the query-builder to eager-load the nodes that are connected to
// the "subject" edge. The optional arguments are used to configure the query builder of the edge.
func (mq *ManifestQuery) WithSubject(opts ...func(*ManifestQuery)) *ManifestQuery {
	query := (&ManifestClient{config: mq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mq.withSubject = query
	return mq
}

// WithReferer tells the query-builder to eager-load the nodes that are connected to
// the "referer" edge. The optional arguments are used to configure the query builder of the edge.
func (mq *ManifestQuery) WithReferer(opts ...func(*ManifestQuery)) *ManifestQuery {
	query := (&ManifestClient{config: mq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mq.withReferer = query
	return mq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		MediaType string `json:"media_type,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Manifest.Query().
//		GroupBy(manifest.FieldMediaType).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (mq *ManifestQuery) GroupBy(field string, fields ...string) *ManifestGroupBy {
	mq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ManifestGroupBy{build: mq}
	grbuild.flds = &mq.ctx.Fields
	grbuild.label = manifest.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		MediaType string `json:"media_type,omitempty"`
//	}
//
//	client.Manifest.Query().
//		Select(manifest.FieldMediaType).
//		Scan(ctx, &v)
func (mq *ManifestQuery) Select(fields ...string) *ManifestSelect {
	mq.ctx.Fields = append(mq.ctx.Fields, fields...)
	sbuild := &ManifestSelect{ManifestQuery: mq}
	sbuild.label = manifest.Label
	sbuild.flds, sbuild.scan = &mq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ManifestSelect configured with the given aggregations.
func (mq *ManifestQuery) Aggregate(fns ...AggregateFunc) *ManifestSelect {
	return mq.Select().Aggregate(fns...)
}

func (mq *ManifestQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range mq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, mq); err != nil {
				return err
			}
		}
	}
	for _, f := range mq.ctx.Fields {
		if !manifest.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if mq.path != nil {
		prev, err := mq.path(ctx)
		if err != nil {
			return err
		}
		mq.sql = prev
	}
	return nil
}

func (mq *ManifestQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Manifest, error) {
	var (
		nodes       = []*Manifest{}
		withFKs     = mq.withFKs
		_spec       = mq.querySpec()
		loadedTypes = [4]bool{
			mq.withTags != nil,
			mq.withRepository != nil,
			mq.withSubject != nil,
			mq.withReferer != nil,
		}
	)
	if mq.withRepository != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, manifest.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Manifest).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Manifest{config: mq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := mq.withTags; query != nil {
		if err := mq.loadTags(ctx, query, nodes,
			func(n *Manifest) { n.Edges.Tags = []*ManifestTagReference{} },
			func(n *Manifest, e *ManifestTagReference) { n.Edges.Tags = append(n.Edges.Tags, e) }); err != nil {
			return nil, err
		}
	}
	if query := mq.withRepository; query != nil {
		if err := mq.loadRepository(ctx, query, nodes, nil,
			func(n *Manifest, e *Repository) { n.Edges.Repository = e }); err != nil {
			return nil, err
		}
	}
	if query := mq.withSubject; query != nil {
		if err := mq.loadSubject(ctx, query, nodes,
			func(n *Manifest) { n.Edges.Subject = []*Manifest{} },
			func(n *Manifest, e *Manifest) { n.Edges.Subject = append(n.Edges.Subject, e) }); err != nil {
			return nil, err
		}
	}
	if query := mq.withReferer; query != nil {
		if err := mq.loadReferer(ctx, query, nodes,
			func(n *Manifest) { n.Edges.Referer = []*Manifest{} },
			func(n *Manifest, e *Manifest) { n.Edges.Referer = append(n.Edges.Referer, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (mq *ManifestQuery) loadTags(ctx context.Context, query *ManifestTagReferenceQuery, nodes []*Manifest, init func(*Manifest), assign func(*Manifest, *ManifestTagReference)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Manifest)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.ManifestTagReference(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(manifest.TagsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.manifest_tag_reference_manifests
		if fk == nil {
			return fmt.Errorf(`foreign-key "manifest_tag_reference_manifests" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "manifest_tag_reference_manifests" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (mq *ManifestQuery) loadRepository(ctx context.Context, query *RepositoryQuery, nodes []*Manifest, init func(*Manifest), assign func(*Manifest, *Repository)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Manifest)
	for i := range nodes {
		if nodes[i].repository_manifests == nil {
			continue
		}
		fk := *nodes[i].repository_manifests
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(repository.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "repository_manifests" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (mq *ManifestQuery) loadSubject(ctx context.Context, query *ManifestQuery, nodes []*Manifest, init func(*Manifest), assign func(*Manifest, *Manifest)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*Manifest)
	nids := make(map[int]map[*Manifest]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(manifest.SubjectTable)
		s.Join(joinT).On(s.C(manifest.FieldID), joinT.C(manifest.SubjectPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(manifest.SubjectPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(manifest.SubjectPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*Manifest]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Manifest](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "subject" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (mq *ManifestQuery) loadReferer(ctx context.Context, query *ManifestQuery, nodes []*Manifest, init func(*Manifest), assign func(*Manifest, *Manifest)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*Manifest)
	nids := make(map[int]map[*Manifest]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(manifest.RefererTable)
		s.Join(joinT).On(s.C(manifest.FieldID), joinT.C(manifest.RefererPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(manifest.RefererPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(manifest.RefererPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*Manifest]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Manifest](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "referer" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (mq *ManifestQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mq.querySpec()
	_spec.Node.Columns = mq.ctx.Fields
	if len(mq.ctx.Fields) > 0 {
		_spec.Unique = mq.ctx.Unique != nil && *mq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, mq.driver, _spec)
}

func (mq *ManifestQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(manifest.Table, manifest.Columns, sqlgraph.NewFieldSpec(manifest.FieldID, field.TypeInt))
	_spec.From = mq.sql
	if unique := mq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if mq.path != nil {
		_spec.Unique = true
	}
	if fields := mq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, manifest.FieldID)
		for i := range fields {
			if fields[i] != manifest.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := mq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mq *ManifestQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mq.driver.Dialect())
	t1 := builder.Table(manifest.Table)
	columns := mq.ctx.Fields
	if len(columns) == 0 {
		columns = manifest.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mq.sql != nil {
		selector = mq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mq.ctx.Unique != nil && *mq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range mq.predicates {
		p(selector)
	}
	for _, p := range mq.order {
		p(selector)
	}
	if offset := mq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ManifestGroupBy is the group-by builder for Manifest entities.
type ManifestGroupBy struct {
	selector
	build *ManifestQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mgb *ManifestGroupBy) Aggregate(fns ...AggregateFunc) *ManifestGroupBy {
	mgb.fns = append(mgb.fns, fns...)
	return mgb
}

// Scan applies the selector query and scans the result into the given value.
func (mgb *ManifestGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mgb.build.ctx, ent.OpQueryGroupBy)
	if err := mgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ManifestQuery, *ManifestGroupBy](ctx, mgb.build, mgb, mgb.build.inters, v)
}

func (mgb *ManifestGroupBy) sqlScan(ctx context.Context, root *ManifestQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(mgb.fns))
	for _, fn := range mgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*mgb.flds)+len(mgb.fns))
		for _, f := range *mgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*mgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ManifestSelect is the builder for selecting fields of Manifest entities.
type ManifestSelect struct {
	*ManifestQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ms *ManifestSelect) Aggregate(fns ...AggregateFunc) *ManifestSelect {
	ms.fns = append(ms.fns, fns...)
	return ms
}

// Scan applies the selector query and scans the result into the given value.
func (ms *ManifestSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ms.ctx, ent.OpQuerySelect)
	if err := ms.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ManifestQuery, *ManifestSelect](ctx, ms.ManifestQuery, ms, ms.inters, v)
}

func (ms *ManifestSelect) sqlScan(ctx context.Context, root *ManifestQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ms.fns))
	for _, fn := range ms.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ms.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
