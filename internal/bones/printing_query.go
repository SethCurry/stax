// Code generated by ent, DO NOT EDIT.

package bones

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/SethCurry/stax/internal/bones/artist"
	"github.com/SethCurry/stax/internal/bones/cardface"
	"github.com/SethCurry/stax/internal/bones/predicate"
	"github.com/SethCurry/stax/internal/bones/printing"
	"github.com/SethCurry/stax/internal/bones/printingimage"
	"github.com/SethCurry/stax/internal/bones/set"
)

// PrintingQuery is the builder for querying Printing entities.
type PrintingQuery struct {
	config
	ctx          *QueryContext
	order        []printing.OrderOption
	inters       []Interceptor
	predicates   []predicate.Printing
	withArtist   *ArtistQuery
	withSet      *SetQuery
	withCardFace *CardFaceQuery
	withImages   *PrintingImageQuery
	withFKs      bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PrintingQuery builder.
func (pq *PrintingQuery) Where(ps ...predicate.Printing) *PrintingQuery {
	pq.predicates = append(pq.predicates, ps...)
	return pq
}

// Limit the number of records to be returned by this query.
func (pq *PrintingQuery) Limit(limit int) *PrintingQuery {
	pq.ctx.Limit = &limit
	return pq
}

// Offset to start from.
func (pq *PrintingQuery) Offset(offset int) *PrintingQuery {
	pq.ctx.Offset = &offset
	return pq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pq *PrintingQuery) Unique(unique bool) *PrintingQuery {
	pq.ctx.Unique = &unique
	return pq
}

// Order specifies how the records should be ordered.
func (pq *PrintingQuery) Order(o ...printing.OrderOption) *PrintingQuery {
	pq.order = append(pq.order, o...)
	return pq
}

// QueryArtist chains the current query on the "artist" edge.
func (pq *PrintingQuery) QueryArtist() *ArtistQuery {
	query := (&ArtistClient{config: pq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(printing.Table, printing.FieldID, selector),
			sqlgraph.To(artist.Table, artist.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, printing.ArtistTable, printing.ArtistColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySet chains the current query on the "set" edge.
func (pq *PrintingQuery) QuerySet() *SetQuery {
	query := (&SetClient{config: pq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(printing.Table, printing.FieldID, selector),
			sqlgraph.To(set.Table, set.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, printing.SetTable, printing.SetColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCardFace chains the current query on the "card_face" edge.
func (pq *PrintingQuery) QueryCardFace() *CardFaceQuery {
	query := (&CardFaceClient{config: pq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(printing.Table, printing.FieldID, selector),
			sqlgraph.To(cardface.Table, cardface.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, printing.CardFaceTable, printing.CardFaceColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryImages chains the current query on the "images" edge.
func (pq *PrintingQuery) QueryImages() *PrintingImageQuery {
	query := (&PrintingImageClient{config: pq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(printing.Table, printing.FieldID, selector),
			sqlgraph.To(printingimage.Table, printingimage.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, printing.ImagesTable, printing.ImagesColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Printing entity from the query.
// Returns a *NotFoundError when no Printing was found.
func (pq *PrintingQuery) First(ctx context.Context) (*Printing, error) {
	nodes, err := pq.Limit(1).All(setContextOp(ctx, pq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{printing.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pq *PrintingQuery) FirstX(ctx context.Context) *Printing {
	node, err := pq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Printing ID from the query.
// Returns a *NotFoundError when no Printing ID was found.
func (pq *PrintingQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pq.Limit(1).IDs(setContextOp(ctx, pq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{printing.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pq *PrintingQuery) FirstIDX(ctx context.Context) int {
	id, err := pq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Printing entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Printing entity is found.
// Returns a *NotFoundError when no Printing entities are found.
func (pq *PrintingQuery) Only(ctx context.Context) (*Printing, error) {
	nodes, err := pq.Limit(2).All(setContextOp(ctx, pq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{printing.Label}
	default:
		return nil, &NotSingularError{printing.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pq *PrintingQuery) OnlyX(ctx context.Context) *Printing {
	node, err := pq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Printing ID in the query.
// Returns a *NotSingularError when more than one Printing ID is found.
// Returns a *NotFoundError when no entities are found.
func (pq *PrintingQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pq.Limit(2).IDs(setContextOp(ctx, pq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{printing.Label}
	default:
		err = &NotSingularError{printing.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pq *PrintingQuery) OnlyIDX(ctx context.Context) int {
	id, err := pq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Printings.
func (pq *PrintingQuery) All(ctx context.Context) ([]*Printing, error) {
	ctx = setContextOp(ctx, pq.ctx, "All")
	if err := pq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Printing, *PrintingQuery]()
	return withInterceptors[[]*Printing](ctx, pq, qr, pq.inters)
}

// AllX is like All, but panics if an error occurs.
func (pq *PrintingQuery) AllX(ctx context.Context) []*Printing {
	nodes, err := pq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Printing IDs.
func (pq *PrintingQuery) IDs(ctx context.Context) (ids []int, err error) {
	if pq.ctx.Unique == nil && pq.path != nil {
		pq.Unique(true)
	}
	ctx = setContextOp(ctx, pq.ctx, "IDs")
	if err = pq.Select(printing.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pq *PrintingQuery) IDsX(ctx context.Context) []int {
	ids, err := pq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pq *PrintingQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, pq.ctx, "Count")
	if err := pq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, pq, querierCount[*PrintingQuery](), pq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (pq *PrintingQuery) CountX(ctx context.Context) int {
	count, err := pq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pq *PrintingQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, pq.ctx, "Exist")
	switch _, err := pq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("bones: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (pq *PrintingQuery) ExistX(ctx context.Context) bool {
	exist, err := pq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PrintingQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pq *PrintingQuery) Clone() *PrintingQuery {
	if pq == nil {
		return nil
	}
	return &PrintingQuery{
		config:       pq.config,
		ctx:          pq.ctx.Clone(),
		order:        append([]printing.OrderOption{}, pq.order...),
		inters:       append([]Interceptor{}, pq.inters...),
		predicates:   append([]predicate.Printing{}, pq.predicates...),
		withArtist:   pq.withArtist.Clone(),
		withSet:      pq.withSet.Clone(),
		withCardFace: pq.withCardFace.Clone(),
		withImages:   pq.withImages.Clone(),
		// clone intermediate query.
		sql:  pq.sql.Clone(),
		path: pq.path,
	}
}

// WithArtist tells the query-builder to eager-load the nodes that are connected to
// the "artist" edge. The optional arguments are used to configure the query builder of the edge.
func (pq *PrintingQuery) WithArtist(opts ...func(*ArtistQuery)) *PrintingQuery {
	query := (&ArtistClient{config: pq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pq.withArtist = query
	return pq
}

// WithSet tells the query-builder to eager-load the nodes that are connected to
// the "set" edge. The optional arguments are used to configure the query builder of the edge.
func (pq *PrintingQuery) WithSet(opts ...func(*SetQuery)) *PrintingQuery {
	query := (&SetClient{config: pq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pq.withSet = query
	return pq
}

// WithCardFace tells the query-builder to eager-load the nodes that are connected to
// the "card_face" edge. The optional arguments are used to configure the query builder of the edge.
func (pq *PrintingQuery) WithCardFace(opts ...func(*CardFaceQuery)) *PrintingQuery {
	query := (&CardFaceClient{config: pq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pq.withCardFace = query
	return pq
}

// WithImages tells the query-builder to eager-load the nodes that are connected to
// the "images" edge. The optional arguments are used to configure the query builder of the edge.
func (pq *PrintingQuery) WithImages(opts ...func(*PrintingImageQuery)) *PrintingQuery {
	query := (&PrintingImageClient{config: pq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pq.withImages = query
	return pq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Rarity printing.Rarity `json:"rarity,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Printing.Query().
//		GroupBy(printing.FieldRarity).
//		Aggregate(bones.Count()).
//		Scan(ctx, &v)
func (pq *PrintingQuery) GroupBy(field string, fields ...string) *PrintingGroupBy {
	pq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PrintingGroupBy{build: pq}
	grbuild.flds = &pq.ctx.Fields
	grbuild.label = printing.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Rarity printing.Rarity `json:"rarity,omitempty"`
//	}
//
//	client.Printing.Query().
//		Select(printing.FieldRarity).
//		Scan(ctx, &v)
func (pq *PrintingQuery) Select(fields ...string) *PrintingSelect {
	pq.ctx.Fields = append(pq.ctx.Fields, fields...)
	sbuild := &PrintingSelect{PrintingQuery: pq}
	sbuild.label = printing.Label
	sbuild.flds, sbuild.scan = &pq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PrintingSelect configured with the given aggregations.
func (pq *PrintingQuery) Aggregate(fns ...AggregateFunc) *PrintingSelect {
	return pq.Select().Aggregate(fns...)
}

func (pq *PrintingQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range pq.inters {
		if inter == nil {
			return fmt.Errorf("bones: uninitialized interceptor (forgotten import bones/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, pq); err != nil {
				return err
			}
		}
	}
	for _, f := range pq.ctx.Fields {
		if !printing.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("bones: invalid field %q for query", f)}
		}
	}
	if pq.path != nil {
		prev, err := pq.path(ctx)
		if err != nil {
			return err
		}
		pq.sql = prev
	}
	return nil
}

func (pq *PrintingQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Printing, error) {
	var (
		nodes       = []*Printing{}
		withFKs     = pq.withFKs
		_spec       = pq.querySpec()
		loadedTypes = [4]bool{
			pq.withArtist != nil,
			pq.withSet != nil,
			pq.withCardFace != nil,
			pq.withImages != nil,
		}
	)
	if pq.withArtist != nil || pq.withSet != nil || pq.withCardFace != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, printing.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Printing).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Printing{config: pq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, pq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := pq.withArtist; query != nil {
		if err := pq.loadArtist(ctx, query, nodes, nil,
			func(n *Printing, e *Artist) { n.Edges.Artist = e }); err != nil {
			return nil, err
		}
	}
	if query := pq.withSet; query != nil {
		if err := pq.loadSet(ctx, query, nodes, nil,
			func(n *Printing, e *Set) { n.Edges.Set = e }); err != nil {
			return nil, err
		}
	}
	if query := pq.withCardFace; query != nil {
		if err := pq.loadCardFace(ctx, query, nodes, nil,
			func(n *Printing, e *CardFace) { n.Edges.CardFace = e }); err != nil {
			return nil, err
		}
	}
	if query := pq.withImages; query != nil {
		if err := pq.loadImages(ctx, query, nodes,
			func(n *Printing) { n.Edges.Images = []*PrintingImage{} },
			func(n *Printing, e *PrintingImage) { n.Edges.Images = append(n.Edges.Images, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (pq *PrintingQuery) loadArtist(ctx context.Context, query *ArtistQuery, nodes []*Printing, init func(*Printing), assign func(*Printing, *Artist)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Printing)
	for i := range nodes {
		if nodes[i].printing_artist == nil {
			continue
		}
		fk := *nodes[i].printing_artist
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(artist.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "printing_artist" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (pq *PrintingQuery) loadSet(ctx context.Context, query *SetQuery, nodes []*Printing, init func(*Printing), assign func(*Printing, *Set)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Printing)
	for i := range nodes {
		if nodes[i].printing_set == nil {
			continue
		}
		fk := *nodes[i].printing_set
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(set.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "printing_set" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (pq *PrintingQuery) loadCardFace(ctx context.Context, query *CardFaceQuery, nodes []*Printing, init func(*Printing), assign func(*Printing, *CardFace)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Printing)
	for i := range nodes {
		if nodes[i].printing_card_face == nil {
			continue
		}
		fk := *nodes[i].printing_card_face
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(cardface.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "printing_card_face" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (pq *PrintingQuery) loadImages(ctx context.Context, query *PrintingImageQuery, nodes []*Printing, init func(*Printing), assign func(*Printing, *PrintingImage)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Printing)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.PrintingImage(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(printing.ImagesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.printing_image_printing
		if fk == nil {
			return fmt.Errorf(`foreign-key "printing_image_printing" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "printing_image_printing" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (pq *PrintingQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pq.querySpec()
	_spec.Node.Columns = pq.ctx.Fields
	if len(pq.ctx.Fields) > 0 {
		_spec.Unique = pq.ctx.Unique != nil && *pq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, pq.driver, _spec)
}

func (pq *PrintingQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(printing.Table, printing.Columns, sqlgraph.NewFieldSpec(printing.FieldID, field.TypeInt))
	_spec.From = pq.sql
	if unique := pq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if pq.path != nil {
		_spec.Unique = true
	}
	if fields := pq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, printing.FieldID)
		for i := range fields {
			if fields[i] != printing.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := pq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pq *PrintingQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pq.driver.Dialect())
	t1 := builder.Table(printing.Table)
	columns := pq.ctx.Fields
	if len(columns) == 0 {
		columns = printing.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pq.sql != nil {
		selector = pq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pq.ctx.Unique != nil && *pq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range pq.predicates {
		p(selector)
	}
	for _, p := range pq.order {
		p(selector)
	}
	if offset := pq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// PrintingGroupBy is the group-by builder for Printing entities.
type PrintingGroupBy struct {
	selector
	build *PrintingQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pgb *PrintingGroupBy) Aggregate(fns ...AggregateFunc) *PrintingGroupBy {
	pgb.fns = append(pgb.fns, fns...)
	return pgb
}

// Scan applies the selector query and scans the result into the given value.
func (pgb *PrintingGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pgb.build.ctx, "GroupBy")
	if err := pgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PrintingQuery, *PrintingGroupBy](ctx, pgb.build, pgb, pgb.build.inters, v)
}

func (pgb *PrintingGroupBy) sqlScan(ctx context.Context, root *PrintingQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(pgb.fns))
	for _, fn := range pgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*pgb.flds)+len(pgb.fns))
		for _, f := range *pgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*pgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PrintingSelect is the builder for selecting fields of Printing entities.
type PrintingSelect struct {
	*PrintingQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ps *PrintingSelect) Aggregate(fns ...AggregateFunc) *PrintingSelect {
	ps.fns = append(ps.fns, fns...)
	return ps
}

// Scan applies the selector query and scans the result into the given value.
func (ps *PrintingSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ps.ctx, "Select")
	if err := ps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PrintingQuery, *PrintingSelect](ctx, ps.PrintingQuery, ps, ps.inters, v)
}

func (ps *PrintingSelect) sqlScan(ctx context.Context, root *PrintingQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ps.fns))
	for _, fn := range ps.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ps.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
