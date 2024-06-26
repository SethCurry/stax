// Code generated by ent, DO NOT EDIT.

package bones

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/SethCurry/stax/internal/bones/cardface"
	"github.com/SethCurry/stax/internal/bones/predicate"
)

// CardFaceDelete is the builder for deleting a CardFace entity.
type CardFaceDelete struct {
	config
	hooks    []Hook
	mutation *CardFaceMutation
}

// Where appends a list predicates to the CardFaceDelete builder.
func (cfd *CardFaceDelete) Where(ps ...predicate.CardFace) *CardFaceDelete {
	cfd.mutation.Where(ps...)
	return cfd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cfd *CardFaceDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cfd.sqlExec, cfd.mutation, cfd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cfd *CardFaceDelete) ExecX(ctx context.Context) int {
	n, err := cfd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cfd *CardFaceDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(cardface.Table, sqlgraph.NewFieldSpec(cardface.FieldID, field.TypeInt))
	if ps := cfd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cfd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cfd.mutation.done = true
	return affected, err
}

// CardFaceDeleteOne is the builder for deleting a single CardFace entity.
type CardFaceDeleteOne struct {
	cfd *CardFaceDelete
}

// Where appends a list predicates to the CardFaceDelete builder.
func (cfdo *CardFaceDeleteOne) Where(ps ...predicate.CardFace) *CardFaceDeleteOne {
	cfdo.cfd.mutation.Where(ps...)
	return cfdo
}

// Exec executes the deletion query.
func (cfdo *CardFaceDeleteOne) Exec(ctx context.Context) error {
	n, err := cfdo.cfd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{cardface.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cfdo *CardFaceDeleteOne) ExecX(ctx context.Context) {
	if err := cfdo.Exec(ctx); err != nil {
		panic(err)
	}
}
