// Code generated by ent, DO NOT EDIT.

package bones

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/SethCurry/stax/internal/bones/artist"
	"github.com/SethCurry/stax/internal/bones/cardface"
	"github.com/SethCurry/stax/internal/bones/printing"
	"github.com/SethCurry/stax/internal/bones/printingimage"
	"github.com/SethCurry/stax/internal/bones/set"
)

// PrintingCreate is the builder for creating a Printing entity.
type PrintingCreate struct {
	config
	mutation *PrintingMutation
	hooks    []Hook
}

// SetRarity sets the "rarity" field.
func (pc *PrintingCreate) SetRarity(pr printing.Rarity) *PrintingCreate {
	pc.mutation.SetRarity(pr)
	return pc
}

// SetArtistID sets the "artist" edge to the Artist entity by ID.
func (pc *PrintingCreate) SetArtistID(id int) *PrintingCreate {
	pc.mutation.SetArtistID(id)
	return pc
}

// SetNillableArtistID sets the "artist" edge to the Artist entity by ID if the given value is not nil.
func (pc *PrintingCreate) SetNillableArtistID(id *int) *PrintingCreate {
	if id != nil {
		pc = pc.SetArtistID(*id)
	}
	return pc
}

// SetArtist sets the "artist" edge to the Artist entity.
func (pc *PrintingCreate) SetArtist(a *Artist) *PrintingCreate {
	return pc.SetArtistID(a.ID)
}

// SetSetID sets the "set" edge to the Set entity by ID.
func (pc *PrintingCreate) SetSetID(id int) *PrintingCreate {
	pc.mutation.SetSetID(id)
	return pc
}

// SetNillableSetID sets the "set" edge to the Set entity by ID if the given value is not nil.
func (pc *PrintingCreate) SetNillableSetID(id *int) *PrintingCreate {
	if id != nil {
		pc = pc.SetSetID(*id)
	}
	return pc
}

// SetSet sets the "set" edge to the Set entity.
func (pc *PrintingCreate) SetSet(s *Set) *PrintingCreate {
	return pc.SetSetID(s.ID)
}

// SetCardFaceID sets the "card_face" edge to the CardFace entity by ID.
func (pc *PrintingCreate) SetCardFaceID(id int) *PrintingCreate {
	pc.mutation.SetCardFaceID(id)
	return pc
}

// SetNillableCardFaceID sets the "card_face" edge to the CardFace entity by ID if the given value is not nil.
func (pc *PrintingCreate) SetNillableCardFaceID(id *int) *PrintingCreate {
	if id != nil {
		pc = pc.SetCardFaceID(*id)
	}
	return pc
}

// SetCardFace sets the "card_face" edge to the CardFace entity.
func (pc *PrintingCreate) SetCardFace(c *CardFace) *PrintingCreate {
	return pc.SetCardFaceID(c.ID)
}

// AddImageIDs adds the "images" edge to the PrintingImage entity by IDs.
func (pc *PrintingCreate) AddImageIDs(ids ...int) *PrintingCreate {
	pc.mutation.AddImageIDs(ids...)
	return pc
}

// AddImages adds the "images" edges to the PrintingImage entity.
func (pc *PrintingCreate) AddImages(p ...*PrintingImage) *PrintingCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pc.AddImageIDs(ids...)
}

// Mutation returns the PrintingMutation object of the builder.
func (pc *PrintingCreate) Mutation() *PrintingMutation {
	return pc.mutation
}

// Save creates the Printing in the database.
func (pc *PrintingCreate) Save(ctx context.Context) (*Printing, error) {
	return withHooks(ctx, pc.sqlSave, pc.mutation, pc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *PrintingCreate) SaveX(ctx context.Context) *Printing {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pc *PrintingCreate) Exec(ctx context.Context) error {
	_, err := pc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pc *PrintingCreate) ExecX(ctx context.Context) {
	if err := pc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pc *PrintingCreate) check() error {
	if _, ok := pc.mutation.Rarity(); !ok {
		return &ValidationError{Name: "rarity", err: errors.New(`bones: missing required field "Printing.rarity"`)}
	}
	if v, ok := pc.mutation.Rarity(); ok {
		if err := printing.RarityValidator(v); err != nil {
			return &ValidationError{Name: "rarity", err: fmt.Errorf(`bones: validator failed for field "Printing.rarity": %w`, err)}
		}
	}
	return nil
}

func (pc *PrintingCreate) sqlSave(ctx context.Context) (*Printing, error) {
	if err := pc.check(); err != nil {
		return nil, err
	}
	_node, _spec := pc.createSpec()
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	pc.mutation.id = &_node.ID
	pc.mutation.done = true
	return _node, nil
}

func (pc *PrintingCreate) createSpec() (*Printing, *sqlgraph.CreateSpec) {
	var (
		_node = &Printing{config: pc.config}
		_spec = sqlgraph.NewCreateSpec(printing.Table, sqlgraph.NewFieldSpec(printing.FieldID, field.TypeInt))
	)
	if value, ok := pc.mutation.Rarity(); ok {
		_spec.SetField(printing.FieldRarity, field.TypeEnum, value)
		_node.Rarity = value
	}
	if nodes := pc.mutation.ArtistIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   printing.ArtistTable,
			Columns: []string{printing.ArtistColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(artist.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.printing_artist = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.SetIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   printing.SetTable,
			Columns: []string{printing.SetColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(set.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.printing_set = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.CardFaceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   printing.CardFaceTable,
			Columns: []string{printing.CardFaceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cardface.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.printing_card_face = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.ImagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   printing.ImagesTable,
			Columns: []string{printing.ImagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(printingimage.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// PrintingCreateBulk is the builder for creating many Printing entities in bulk.
type PrintingCreateBulk struct {
	config
	err      error
	builders []*PrintingCreate
}

// Save creates the Printing entities in the database.
func (pcb *PrintingCreateBulk) Save(ctx context.Context) ([]*Printing, error) {
	if pcb.err != nil {
		return nil, pcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(pcb.builders))
	nodes := make([]*Printing, len(pcb.builders))
	mutators := make([]Mutator, len(pcb.builders))
	for i := range pcb.builders {
		func(i int, root context.Context) {
			builder := pcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PrintingMutation)
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
					_, err = mutators[i+1].Mutate(root, pcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, pcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pcb *PrintingCreateBulk) SaveX(ctx context.Context) []*Printing {
	v, err := pcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pcb *PrintingCreateBulk) Exec(ctx context.Context) error {
	_, err := pcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pcb *PrintingCreateBulk) ExecX(ctx context.Context) {
	if err := pcb.Exec(ctx); err != nil {
		panic(err)
	}
}
