// Code generated by ent, DO NOT EDIT.

package repository

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Repository {
	return predicate.Repository(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Repository {
	return predicate.Repository(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Repository {
	return predicate.Repository(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Repository {
	return predicate.Repository(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Repository {
	return predicate.Repository(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Repository {
	return predicate.Repository(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Repository {
	return predicate.Repository(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Repository {
	return predicate.Repository(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Repository {
	return predicate.Repository(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Repository {
	return predicate.Repository(sql.FieldEQ(FieldName, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Repository {
	return predicate.Repository(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Repository {
	return predicate.Repository(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Repository {
	return predicate.Repository(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Repository {
	return predicate.Repository(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Repository {
	return predicate.Repository(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Repository {
	return predicate.Repository(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Repository {
	return predicate.Repository(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Repository {
	return predicate.Repository(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Repository {
	return predicate.Repository(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Repository {
	return predicate.Repository(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Repository {
	return predicate.Repository(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Repository {
	return predicate.Repository(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Repository {
	return predicate.Repository(sql.FieldContainsFold(FieldName, v))
}

// HasManifests applies the HasEdge predicate on the "manifests" edge.
func HasManifests() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ManifestsTable, ManifestsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasManifestsWith applies the HasEdge predicate on the "manifests" edge with a given conditions (other predicates).
func HasManifestsWith(preds ...predicate.Manifest) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := newManifestsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRegistry applies the HasEdge predicate on the "registry" edge.
func HasRegistry() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, RegistryTable, RegistryColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRegistryWith applies the HasEdge predicate on the "registry" edge with a given conditions (other predicates).
func HasRegistryWith(preds ...predicate.Registry) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := newRegistryStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Repository) predicate.Repository {
	return predicate.Repository(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Repository) predicate.Repository {
	return predicate.Repository(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Repository) predicate.Repository {
	return predicate.Repository(sql.NotPredicates(p))
}
