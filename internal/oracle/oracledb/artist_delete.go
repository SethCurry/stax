// Code generated by ent, DO NOT EDIT.

package oracledb

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/SethCurry/stax/internal/oracle/oracledb/artist"
	"github.com/SethCurry/stax/internal/oracle/oracledb/predicate"
)

// ArtistDelete is the builder for deleting a Artist entity.
type ArtistDelete struct {
	config
	hooks    []Hook
	mutation *ArtistMutation
}

// Where appends a list predicates to the ArtistDelete builder.
func (ad *ArtistDelete) Where(ps ...predicate.Artist) *ArtistDelete {
	ad.mutation.Where(ps...)
	return ad
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ad *ArtistDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ad.sqlExec, ad.mutation, ad.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ad *ArtistDelete) ExecX(ctx context.Context) int {
	n, err := ad.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ad *ArtistDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(artist.Table, sqlgraph.NewFieldSpec(artist.FieldID, field.TypeInt))
	if ps := ad.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ad.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ad.mutation.done = true
	return affected, err
}

// ArtistDeleteOne is the builder for deleting a single Artist entity.
type ArtistDeleteOne struct {
	ad *ArtistDelete
}

// Where appends a list predicates to the ArtistDelete builder.
func (ado *ArtistDeleteOne) Where(ps ...predicate.Artist) *ArtistDeleteOne {
	ado.ad.mutation.Where(ps...)
	return ado
}

// Exec executes the deletion query.
func (ado *ArtistDeleteOne) Exec(ctx context.Context) error {
	n, err := ado.ad.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{artist.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ado *ArtistDeleteOne) ExecX(ctx context.Context) {
	if err := ado.Exec(ctx); err != nil {
		panic(err)
	}
}
