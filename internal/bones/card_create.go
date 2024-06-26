// Code generated by ent, DO NOT EDIT.

package bones

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/SethCurry/stax/internal/bones/card"
	"github.com/SethCurry/stax/internal/bones/cardface"
	"github.com/SethCurry/stax/internal/bones/ruling"
)

// CardCreate is the builder for creating a Card entity.
type CardCreate struct {
	config
	mutation *CardMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (cc *CardCreate) SetName(s string) *CardCreate {
	cc.mutation.SetName(s)
	return cc
}

// SetOracleID sets the "oracle_id" field.
func (cc *CardCreate) SetOracleID(s string) *CardCreate {
	cc.mutation.SetOracleID(s)
	return cc
}

// SetColorIdentity sets the "color_identity" field.
func (cc *CardCreate) SetColorIdentity(u uint8) *CardCreate {
	cc.mutation.SetColorIdentity(u)
	return cc
}

// AddFaceIDs adds the "faces" edge to the CardFace entity by IDs.
func (cc *CardCreate) AddFaceIDs(ids ...int) *CardCreate {
	cc.mutation.AddFaceIDs(ids...)
	return cc
}

// AddFaces adds the "faces" edges to the CardFace entity.
func (cc *CardCreate) AddFaces(c ...*CardFace) *CardCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cc.AddFaceIDs(ids...)
}

// AddRulingIDs adds the "rulings" edge to the Ruling entity by IDs.
func (cc *CardCreate) AddRulingIDs(ids ...int) *CardCreate {
	cc.mutation.AddRulingIDs(ids...)
	return cc
}

// AddRulings adds the "rulings" edges to the Ruling entity.
func (cc *CardCreate) AddRulings(r ...*Ruling) *CardCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return cc.AddRulingIDs(ids...)
}

// Mutation returns the CardMutation object of the builder.
func (cc *CardCreate) Mutation() *CardMutation {
	return cc.mutation
}

// Save creates the Card in the database.
func (cc *CardCreate) Save(ctx context.Context) (*Card, error) {
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CardCreate) SaveX(ctx context.Context) *Card {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *CardCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *CardCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *CardCreate) check() error {
	if _, ok := cc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`bones: missing required field "Card.name"`)}
	}
	if v, ok := cc.mutation.Name(); ok {
		if err := card.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`bones: validator failed for field "Card.name": %w`, err)}
		}
	}
	if _, ok := cc.mutation.OracleID(); !ok {
		return &ValidationError{Name: "oracle_id", err: errors.New(`bones: missing required field "Card.oracle_id"`)}
	}
	if v, ok := cc.mutation.OracleID(); ok {
		if err := card.OracleIDValidator(v); err != nil {
			return &ValidationError{Name: "oracle_id", err: fmt.Errorf(`bones: validator failed for field "Card.oracle_id": %w`, err)}
		}
	}
	if _, ok := cc.mutation.ColorIdentity(); !ok {
		return &ValidationError{Name: "color_identity", err: errors.New(`bones: missing required field "Card.color_identity"`)}
	}
	return nil
}

func (cc *CardCreate) sqlSave(ctx context.Context) (*Card, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *CardCreate) createSpec() (*Card, *sqlgraph.CreateSpec) {
	var (
		_node = &Card{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(card.Table, sqlgraph.NewFieldSpec(card.FieldID, field.TypeInt))
	)
	if value, ok := cc.mutation.Name(); ok {
		_spec.SetField(card.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := cc.mutation.OracleID(); ok {
		_spec.SetField(card.FieldOracleID, field.TypeString, value)
		_node.OracleID = value
	}
	if value, ok := cc.mutation.ColorIdentity(); ok {
		_spec.SetField(card.FieldColorIdentity, field.TypeUint8, value)
		_node.ColorIdentity = value
	}
	if nodes := cc.mutation.FacesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   card.FacesTable,
			Columns: []string{card.FacesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cardface.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.RulingsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   card.RulingsTable,
			Columns: []string{card.RulingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(ruling.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// CardCreateBulk is the builder for creating many Card entities in bulk.
type CardCreateBulk struct {
	config
	err      error
	builders []*CardCreate
}

// Save creates the Card entities in the database.
func (ccb *CardCreateBulk) Save(ctx context.Context) ([]*Card, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Card, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CardMutation)
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
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *CardCreateBulk) SaveX(ctx context.Context) []*Card {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *CardCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *CardCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
