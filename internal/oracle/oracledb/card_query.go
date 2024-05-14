// Code generated by ent, DO NOT EDIT.

package oracledb

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/SethCurry/stax/internal/oracle/oracledb/card"
	"github.com/SethCurry/stax/internal/oracle/oracledb/cardface"
	"github.com/SethCurry/stax/internal/oracle/oracledb/predicate"
	"github.com/SethCurry/stax/internal/oracle/oracledb/ruling"
)

// CardQuery is the builder for querying Card entities.
type CardQuery struct {
	config
	ctx         *QueryContext
	order       []card.OrderOption
	inters      []Interceptor
	predicates  []predicate.Card
	withFaces   *CardFaceQuery
	withRulings *RulingQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CardQuery builder.
func (cq *CardQuery) Where(ps ...predicate.Card) *CardQuery {
	cq.predicates = append(cq.predicates, ps...)
	return cq
}

// Limit the number of records to be returned by this query.
func (cq *CardQuery) Limit(limit int) *CardQuery {
	cq.ctx.Limit = &limit
	return cq
}

// Offset to start from.
func (cq *CardQuery) Offset(offset int) *CardQuery {
	cq.ctx.Offset = &offset
	return cq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (cq *CardQuery) Unique(unique bool) *CardQuery {
	cq.ctx.Unique = &unique
	return cq
}

// Order specifies how the records should be ordered.
func (cq *CardQuery) Order(o ...card.OrderOption) *CardQuery {
	cq.order = append(cq.order, o...)
	return cq
}

// QueryFaces chains the current query on the "faces" edge.
func (cq *CardQuery) QueryFaces() *CardFaceQuery {
	query := (&CardFaceClient{config: cq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(card.Table, card.FieldID, selector),
			sqlgraph.To(cardface.Table, cardface.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, card.FacesTable, card.FacesColumn),
		)
		fromU = sqlgraph.SetNeighbors(cq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRulings chains the current query on the "rulings" edge.
func (cq *CardQuery) QueryRulings() *RulingQuery {
	query := (&RulingClient{config: cq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(card.Table, card.FieldID, selector),
			sqlgraph.To(ruling.Table, ruling.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, card.RulingsTable, card.RulingsColumn),
		)
		fromU = sqlgraph.SetNeighbors(cq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Card entity from the query.
// Returns a *NotFoundError when no Card was found.
func (cq *CardQuery) First(ctx context.Context) (*Card, error) {
	nodes, err := cq.Limit(1).All(setContextOp(ctx, cq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{card.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (cq *CardQuery) FirstX(ctx context.Context) *Card {
	node, err := cq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Card ID from the query.
// Returns a *NotFoundError when no Card ID was found.
func (cq *CardQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cq.Limit(1).IDs(setContextOp(ctx, cq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{card.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (cq *CardQuery) FirstIDX(ctx context.Context) int {
	id, err := cq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Card entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Card entity is found.
// Returns a *NotFoundError when no Card entities are found.
func (cq *CardQuery) Only(ctx context.Context) (*Card, error) {
	nodes, err := cq.Limit(2).All(setContextOp(ctx, cq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{card.Label}
	default:
		return nil, &NotSingularError{card.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (cq *CardQuery) OnlyX(ctx context.Context) *Card {
	node, err := cq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Card ID in the query.
// Returns a *NotSingularError when more than one Card ID is found.
// Returns a *NotFoundError when no entities are found.
func (cq *CardQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cq.Limit(2).IDs(setContextOp(ctx, cq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{card.Label}
	default:
		err = &NotSingularError{card.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (cq *CardQuery) OnlyIDX(ctx context.Context) int {
	id, err := cq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Cards.
func (cq *CardQuery) All(ctx context.Context) ([]*Card, error) {
	ctx = setContextOp(ctx, cq.ctx, "All")
	if err := cq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Card, *CardQuery]()
	return withInterceptors[[]*Card](ctx, cq, qr, cq.inters)
}

// AllX is like All, but panics if an error occurs.
func (cq *CardQuery) AllX(ctx context.Context) []*Card {
	nodes, err := cq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Card IDs.
func (cq *CardQuery) IDs(ctx context.Context) (ids []int, err error) {
	if cq.ctx.Unique == nil && cq.path != nil {
		cq.Unique(true)
	}
	ctx = setContextOp(ctx, cq.ctx, "IDs")
	if err = cq.Select(card.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (cq *CardQuery) IDsX(ctx context.Context) []int {
	ids, err := cq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (cq *CardQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, cq.ctx, "Count")
	if err := cq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, cq, querierCount[*CardQuery](), cq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (cq *CardQuery) CountX(ctx context.Context) int {
	count, err := cq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (cq *CardQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, cq.ctx, "Exist")
	switch _, err := cq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("oracledb: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (cq *CardQuery) ExistX(ctx context.Context) bool {
	exist, err := cq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CardQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (cq *CardQuery) Clone() *CardQuery {
	if cq == nil {
		return nil
	}
	return &CardQuery{
		config:      cq.config,
		ctx:         cq.ctx.Clone(),
		order:       append([]card.OrderOption{}, cq.order...),
		inters:      append([]Interceptor{}, cq.inters...),
		predicates:  append([]predicate.Card{}, cq.predicates...),
		withFaces:   cq.withFaces.Clone(),
		withRulings: cq.withRulings.Clone(),
		// clone intermediate query.
		sql:  cq.sql.Clone(),
		path: cq.path,
	}
}

// WithFaces tells the query-builder to eager-load the nodes that are connected to
// the "faces" edge. The optional arguments are used to configure the query builder of the edge.
func (cq *CardQuery) WithFaces(opts ...func(*CardFaceQuery)) *CardQuery {
	query := (&CardFaceClient{config: cq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	cq.withFaces = query
	return cq
}

// WithRulings tells the query-builder to eager-load the nodes that are connected to
// the "rulings" edge. The optional arguments are used to configure the query builder of the edge.
func (cq *CardQuery) WithRulings(opts ...func(*RulingQuery)) *CardQuery {
	query := (&RulingClient{config: cq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	cq.withRulings = query
	return cq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Card.Query().
//		GroupBy(card.FieldName).
//		Aggregate(oracledb.Count()).
//		Scan(ctx, &v)
func (cq *CardQuery) GroupBy(field string, fields ...string) *CardGroupBy {
	cq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CardGroupBy{build: cq}
	grbuild.flds = &cq.ctx.Fields
	grbuild.label = card.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.Card.Query().
//		Select(card.FieldName).
//		Scan(ctx, &v)
func (cq *CardQuery) Select(fields ...string) *CardSelect {
	cq.ctx.Fields = append(cq.ctx.Fields, fields...)
	sbuild := &CardSelect{CardQuery: cq}
	sbuild.label = card.Label
	sbuild.flds, sbuild.scan = &cq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CardSelect configured with the given aggregations.
func (cq *CardQuery) Aggregate(fns ...AggregateFunc) *CardSelect {
	return cq.Select().Aggregate(fns...)
}

func (cq *CardQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range cq.inters {
		if inter == nil {
			return fmt.Errorf("oracledb: uninitialized interceptor (forgotten import oracledb/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, cq); err != nil {
				return err
			}
		}
	}
	for _, f := range cq.ctx.Fields {
		if !card.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("oracledb: invalid field %q for query", f)}
		}
	}
	if cq.path != nil {
		prev, err := cq.path(ctx)
		if err != nil {
			return err
		}
		cq.sql = prev
	}
	return nil
}

func (cq *CardQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Card, error) {
	var (
		nodes       = []*Card{}
		_spec       = cq.querySpec()
		loadedTypes = [2]bool{
			cq.withFaces != nil,
			cq.withRulings != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Card).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Card{config: cq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, cq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := cq.withFaces; query != nil {
		if err := cq.loadFaces(ctx, query, nodes,
			func(n *Card) { n.Edges.Faces = []*CardFace{} },
			func(n *Card, e *CardFace) { n.Edges.Faces = append(n.Edges.Faces, e) }); err != nil {
			return nil, err
		}
	}
	if query := cq.withRulings; query != nil {
		if err := cq.loadRulings(ctx, query, nodes,
			func(n *Card) { n.Edges.Rulings = []*Ruling{} },
			func(n *Card, e *Ruling) { n.Edges.Rulings = append(n.Edges.Rulings, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (cq *CardQuery) loadFaces(ctx context.Context, query *CardFaceQuery, nodes []*Card, init func(*Card), assign func(*Card, *CardFace)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Card)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.CardFace(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(card.FacesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.card_face_card
		if fk == nil {
			return fmt.Errorf(`foreign-key "card_face_card" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "card_face_card" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (cq *CardQuery) loadRulings(ctx context.Context, query *RulingQuery, nodes []*Card, init func(*Card), assign func(*Card, *Ruling)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Card)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Ruling(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(card.RulingsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.ruling_card
		if fk == nil {
			return fmt.Errorf(`foreign-key "ruling_card" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "ruling_card" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (cq *CardQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := cq.querySpec()
	_spec.Node.Columns = cq.ctx.Fields
	if len(cq.ctx.Fields) > 0 {
		_spec.Unique = cq.ctx.Unique != nil && *cq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, cq.driver, _spec)
}

func (cq *CardQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(card.Table, card.Columns, sqlgraph.NewFieldSpec(card.FieldID, field.TypeInt))
	_spec.From = cq.sql
	if unique := cq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if cq.path != nil {
		_spec.Unique = true
	}
	if fields := cq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, card.FieldID)
		for i := range fields {
			if fields[i] != card.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := cq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := cq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := cq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := cq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (cq *CardQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(cq.driver.Dialect())
	t1 := builder.Table(card.Table)
	columns := cq.ctx.Fields
	if len(columns) == 0 {
		columns = card.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if cq.sql != nil {
		selector = cq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if cq.ctx.Unique != nil && *cq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range cq.predicates {
		p(selector)
	}
	for _, p := range cq.order {
		p(selector)
	}
	if offset := cq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := cq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CardGroupBy is the group-by builder for Card entities.
type CardGroupBy struct {
	selector
	build *CardQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cgb *CardGroupBy) Aggregate(fns ...AggregateFunc) *CardGroupBy {
	cgb.fns = append(cgb.fns, fns...)
	return cgb
}

// Scan applies the selector query and scans the result into the given value.
func (cgb *CardGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cgb.build.ctx, "GroupBy")
	if err := cgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CardQuery, *CardGroupBy](ctx, cgb.build, cgb, cgb.build.inters, v)
}

func (cgb *CardGroupBy) sqlScan(ctx context.Context, root *CardQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(cgb.fns))
	for _, fn := range cgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*cgb.flds)+len(cgb.fns))
		for _, f := range *cgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*cgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CardSelect is the builder for selecting fields of Card entities.
type CardSelect struct {
	*CardQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (cs *CardSelect) Aggregate(fns ...AggregateFunc) *CardSelect {
	cs.fns = append(cs.fns, fns...)
	return cs
}

// Scan applies the selector query and scans the result into the given value.
func (cs *CardSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cs.ctx, "Select")
	if err := cs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CardQuery, *CardSelect](ctx, cs.CardQuery, cs, cs.inters, v)
}

func (cs *CardSelect) sqlScan(ctx context.Context, root *CardQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(cs.fns))
	for _, fn := range cs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*cs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
