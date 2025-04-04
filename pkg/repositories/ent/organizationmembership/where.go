// Code generated by ent, DO NOT EDIT.

package organizationmembership

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
)

// UserID applies equality check predicate on the "user_id" field. It's identical to UserIDEQ.
func UserID(v int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldEQ(FieldUserID, v))
}

// OrganizationID applies equality check predicate on the "organization_id" field. It's identical to OrganizationIDEQ.
func OrganizationID(v int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldEQ(FieldOrganizationID, v))
}

// UserIDEQ applies the EQ predicate on the "user_id" field.
func UserIDEQ(v int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldEQ(FieldUserID, v))
}

// UserIDNEQ applies the NEQ predicate on the "user_id" field.
func UserIDNEQ(v int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldNEQ(FieldUserID, v))
}

// UserIDIn applies the In predicate on the "user_id" field.
func UserIDIn(vs ...int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldIn(FieldUserID, vs...))
}

// UserIDNotIn applies the NotIn predicate on the "user_id" field.
func UserIDNotIn(vs ...int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldNotIn(FieldUserID, vs...))
}

// OrganizationIDEQ applies the EQ predicate on the "organization_id" field.
func OrganizationIDEQ(v int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldEQ(FieldOrganizationID, v))
}

// OrganizationIDNEQ applies the NEQ predicate on the "organization_id" field.
func OrganizationIDNEQ(v int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldNEQ(FieldOrganizationID, v))
}

// OrganizationIDIn applies the In predicate on the "organization_id" field.
func OrganizationIDIn(vs ...int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldIn(FieldOrganizationID, vs...))
}

// OrganizationIDNotIn applies the NotIn predicate on the "organization_id" field.
func OrganizationIDNotIn(vs ...int) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldNotIn(FieldOrganizationID, vs...))
}

// RoleEQ applies the EQ predicate on the "role" field.
func RoleEQ(v Role) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldEQ(FieldRole, v))
}

// RoleNEQ applies the NEQ predicate on the "role" field.
func RoleNEQ(v Role) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldNEQ(FieldRole, v))
}

// RoleIn applies the In predicate on the "role" field.
func RoleIn(vs ...Role) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldIn(FieldRole, vs...))
}

// RoleNotIn applies the NotIn predicate on the "role" field.
func RoleNotIn(vs ...Role) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.FieldNotIn(FieldRole, vs...))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.OrganizationMembership {
	return predicate.OrganizationMembership(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, UserColumn),
			sqlgraph.Edge(sqlgraph.M2O, false, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasOrganization applies the HasEdge predicate on the "organization" edge.
func HasOrganization() predicate.OrganizationMembership {
	return predicate.OrganizationMembership(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, OrganizationColumn),
			sqlgraph.Edge(sqlgraph.M2O, false, OrganizationTable, OrganizationColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOrganizationWith applies the HasEdge predicate on the "organization" edge with a given conditions (other predicates).
func HasOrganizationWith(preds ...predicate.Organization) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(func(s *sql.Selector) {
		step := newOrganizationStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.OrganizationMembership) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.OrganizationMembership) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.OrganizationMembership) predicate.OrganizationMembership {
	return predicate.OrganizationMembership(sql.NotPredicates(p))
}
