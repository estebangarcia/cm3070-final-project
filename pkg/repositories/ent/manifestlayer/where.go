// Code generated by ent, DO NOT EDIT.

package manifestlayer

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldLTE(FieldID, id))
}

// MediaType applies equality check predicate on the "media_type" field. It's identical to MediaTypeEQ.
func MediaType(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEQ(FieldMediaType, v))
}

// Digest applies equality check predicate on the "digest" field. It's identical to DigestEQ.
func Digest(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEQ(FieldDigest, v))
}

// Size applies equality check predicate on the "size" field. It's identical to SizeEQ.
func Size(v int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEQ(FieldSize, v))
}

// MediaTypeEQ applies the EQ predicate on the "media_type" field.
func MediaTypeEQ(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEQ(FieldMediaType, v))
}

// MediaTypeNEQ applies the NEQ predicate on the "media_type" field.
func MediaTypeNEQ(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldNEQ(FieldMediaType, v))
}

// MediaTypeIn applies the In predicate on the "media_type" field.
func MediaTypeIn(vs ...string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldIn(FieldMediaType, vs...))
}

// MediaTypeNotIn applies the NotIn predicate on the "media_type" field.
func MediaTypeNotIn(vs ...string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldNotIn(FieldMediaType, vs...))
}

// MediaTypeGT applies the GT predicate on the "media_type" field.
func MediaTypeGT(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldGT(FieldMediaType, v))
}

// MediaTypeGTE applies the GTE predicate on the "media_type" field.
func MediaTypeGTE(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldGTE(FieldMediaType, v))
}

// MediaTypeLT applies the LT predicate on the "media_type" field.
func MediaTypeLT(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldLT(FieldMediaType, v))
}

// MediaTypeLTE applies the LTE predicate on the "media_type" field.
func MediaTypeLTE(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldLTE(FieldMediaType, v))
}

// MediaTypeContains applies the Contains predicate on the "media_type" field.
func MediaTypeContains(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldContains(FieldMediaType, v))
}

// MediaTypeHasPrefix applies the HasPrefix predicate on the "media_type" field.
func MediaTypeHasPrefix(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldHasPrefix(FieldMediaType, v))
}

// MediaTypeHasSuffix applies the HasSuffix predicate on the "media_type" field.
func MediaTypeHasSuffix(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldHasSuffix(FieldMediaType, v))
}

// MediaTypeEqualFold applies the EqualFold predicate on the "media_type" field.
func MediaTypeEqualFold(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEqualFold(FieldMediaType, v))
}

// MediaTypeContainsFold applies the ContainsFold predicate on the "media_type" field.
func MediaTypeContainsFold(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldContainsFold(FieldMediaType, v))
}

// DigestEQ applies the EQ predicate on the "digest" field.
func DigestEQ(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEQ(FieldDigest, v))
}

// DigestNEQ applies the NEQ predicate on the "digest" field.
func DigestNEQ(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldNEQ(FieldDigest, v))
}

// DigestIn applies the In predicate on the "digest" field.
func DigestIn(vs ...string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldIn(FieldDigest, vs...))
}

// DigestNotIn applies the NotIn predicate on the "digest" field.
func DigestNotIn(vs ...string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldNotIn(FieldDigest, vs...))
}

// DigestGT applies the GT predicate on the "digest" field.
func DigestGT(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldGT(FieldDigest, v))
}

// DigestGTE applies the GTE predicate on the "digest" field.
func DigestGTE(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldGTE(FieldDigest, v))
}

// DigestLT applies the LT predicate on the "digest" field.
func DigestLT(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldLT(FieldDigest, v))
}

// DigestLTE applies the LTE predicate on the "digest" field.
func DigestLTE(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldLTE(FieldDigest, v))
}

// DigestContains applies the Contains predicate on the "digest" field.
func DigestContains(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldContains(FieldDigest, v))
}

// DigestHasPrefix applies the HasPrefix predicate on the "digest" field.
func DigestHasPrefix(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldHasPrefix(FieldDigest, v))
}

// DigestHasSuffix applies the HasSuffix predicate on the "digest" field.
func DigestHasSuffix(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldHasSuffix(FieldDigest, v))
}

// DigestEqualFold applies the EqualFold predicate on the "digest" field.
func DigestEqualFold(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEqualFold(FieldDigest, v))
}

// DigestContainsFold applies the ContainsFold predicate on the "digest" field.
func DigestContainsFold(v string) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldContainsFold(FieldDigest, v))
}

// SizeEQ applies the EQ predicate on the "size" field.
func SizeEQ(v int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldEQ(FieldSize, v))
}

// SizeNEQ applies the NEQ predicate on the "size" field.
func SizeNEQ(v int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldNEQ(FieldSize, v))
}

// SizeIn applies the In predicate on the "size" field.
func SizeIn(vs ...int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldIn(FieldSize, vs...))
}

// SizeNotIn applies the NotIn predicate on the "size" field.
func SizeNotIn(vs ...int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldNotIn(FieldSize, vs...))
}

// SizeGT applies the GT predicate on the "size" field.
func SizeGT(v int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldGT(FieldSize, v))
}

// SizeGTE applies the GTE predicate on the "size" field.
func SizeGTE(v int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldGTE(FieldSize, v))
}

// SizeLT applies the LT predicate on the "size" field.
func SizeLT(v int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldLT(FieldSize, v))
}

// SizeLTE applies the LTE predicate on the "size" field.
func SizeLTE(v int32) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.FieldLTE(FieldSize, v))
}

// HasManifest applies the HasEdge predicate on the "manifest" edge.
func HasManifest() predicate.ManifestLayer {
	return predicate.ManifestLayer(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ManifestTable, ManifestColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasManifestWith applies the HasEdge predicate on the "manifest" edge with a given conditions (other predicates).
func HasManifestWith(preds ...predicate.Manifest) predicate.ManifestLayer {
	return predicate.ManifestLayer(func(s *sql.Selector) {
		step := newManifestStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ManifestLayer) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ManifestLayer) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ManifestLayer) predicate.ManifestLayer {
	return predicate.ManifestLayer(sql.NotPredicates(p))
}
