// Code generated by ent, DO NOT EDIT.

package printingimage

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/SethCurry/stax/internal/bones/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldLTE(FieldID, id))
}

// URL applies equality check predicate on the "url" field. It's identical to URLEQ.
func URL(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEQ(FieldURL, v))
}

// LocalPath applies equality check predicate on the "local_path" field. It's identical to LocalPathEQ.
func LocalPath(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEQ(FieldLocalPath, v))
}

// URLEQ applies the EQ predicate on the "url" field.
func URLEQ(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEQ(FieldURL, v))
}

// URLNEQ applies the NEQ predicate on the "url" field.
func URLNEQ(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNEQ(FieldURL, v))
}

// URLIn applies the In predicate on the "url" field.
func URLIn(vs ...string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldIn(FieldURL, vs...))
}

// URLNotIn applies the NotIn predicate on the "url" field.
func URLNotIn(vs ...string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNotIn(FieldURL, vs...))
}

// URLGT applies the GT predicate on the "url" field.
func URLGT(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldGT(FieldURL, v))
}

// URLGTE applies the GTE predicate on the "url" field.
func URLGTE(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldGTE(FieldURL, v))
}

// URLLT applies the LT predicate on the "url" field.
func URLLT(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldLT(FieldURL, v))
}

// URLLTE applies the LTE predicate on the "url" field.
func URLLTE(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldLTE(FieldURL, v))
}

// URLContains applies the Contains predicate on the "url" field.
func URLContains(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldContains(FieldURL, v))
}

// URLHasPrefix applies the HasPrefix predicate on the "url" field.
func URLHasPrefix(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldHasPrefix(FieldURL, v))
}

// URLHasSuffix applies the HasSuffix predicate on the "url" field.
func URLHasSuffix(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldHasSuffix(FieldURL, v))
}

// URLEqualFold applies the EqualFold predicate on the "url" field.
func URLEqualFold(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEqualFold(FieldURL, v))
}

// URLContainsFold applies the ContainsFold predicate on the "url" field.
func URLContainsFold(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldContainsFold(FieldURL, v))
}

// ImageTypeEQ applies the EQ predicate on the "image_type" field.
func ImageTypeEQ(v ImageType) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEQ(FieldImageType, v))
}

// ImageTypeNEQ applies the NEQ predicate on the "image_type" field.
func ImageTypeNEQ(v ImageType) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNEQ(FieldImageType, v))
}

// ImageTypeIn applies the In predicate on the "image_type" field.
func ImageTypeIn(vs ...ImageType) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldIn(FieldImageType, vs...))
}

// ImageTypeNotIn applies the NotIn predicate on the "image_type" field.
func ImageTypeNotIn(vs ...ImageType) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNotIn(FieldImageType, vs...))
}

// LocalPathEQ applies the EQ predicate on the "local_path" field.
func LocalPathEQ(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEQ(FieldLocalPath, v))
}

// LocalPathNEQ applies the NEQ predicate on the "local_path" field.
func LocalPathNEQ(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNEQ(FieldLocalPath, v))
}

// LocalPathIn applies the In predicate on the "local_path" field.
func LocalPathIn(vs ...string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldIn(FieldLocalPath, vs...))
}

// LocalPathNotIn applies the NotIn predicate on the "local_path" field.
func LocalPathNotIn(vs ...string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNotIn(FieldLocalPath, vs...))
}

// LocalPathGT applies the GT predicate on the "local_path" field.
func LocalPathGT(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldGT(FieldLocalPath, v))
}

// LocalPathGTE applies the GTE predicate on the "local_path" field.
func LocalPathGTE(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldGTE(FieldLocalPath, v))
}

// LocalPathLT applies the LT predicate on the "local_path" field.
func LocalPathLT(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldLT(FieldLocalPath, v))
}

// LocalPathLTE applies the LTE predicate on the "local_path" field.
func LocalPathLTE(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldLTE(FieldLocalPath, v))
}

// LocalPathContains applies the Contains predicate on the "local_path" field.
func LocalPathContains(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldContains(FieldLocalPath, v))
}

// LocalPathHasPrefix applies the HasPrefix predicate on the "local_path" field.
func LocalPathHasPrefix(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldHasPrefix(FieldLocalPath, v))
}

// LocalPathHasSuffix applies the HasSuffix predicate on the "local_path" field.
func LocalPathHasSuffix(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldHasSuffix(FieldLocalPath, v))
}

// LocalPathIsNil applies the IsNil predicate on the "local_path" field.
func LocalPathIsNil() predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldIsNull(FieldLocalPath))
}

// LocalPathNotNil applies the NotNil predicate on the "local_path" field.
func LocalPathNotNil() predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldNotNull(FieldLocalPath))
}

// LocalPathEqualFold applies the EqualFold predicate on the "local_path" field.
func LocalPathEqualFold(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldEqualFold(FieldLocalPath, v))
}

// LocalPathContainsFold applies the ContainsFold predicate on the "local_path" field.
func LocalPathContainsFold(v string) predicate.PrintingImage {
	return predicate.PrintingImage(sql.FieldContainsFold(FieldLocalPath, v))
}

// HasPrinting applies the HasEdge predicate on the "printing" edge.
func HasPrinting() predicate.PrintingImage {
	return predicate.PrintingImage(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, PrintingTable, PrintingColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPrintingWith applies the HasEdge predicate on the "printing" edge with a given conditions (other predicates).
func HasPrintingWith(preds ...predicate.Printing) predicate.PrintingImage {
	return predicate.PrintingImage(func(s *sql.Selector) {
		step := newPrintingStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PrintingImage) predicate.PrintingImage {
	return predicate.PrintingImage(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PrintingImage) predicate.PrintingImage {
	return predicate.PrintingImage(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.PrintingImage) predicate.PrintingImage {
	return predicate.PrintingImage(sql.NotPredicates(p))
}
