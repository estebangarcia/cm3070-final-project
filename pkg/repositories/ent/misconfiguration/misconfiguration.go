// Code generated by ent, DO NOT EDIT.

package misconfiguration

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the misconfiguration type in the database.
	Label = "misconfiguration"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldMisconfigurationID holds the string denoting the misconfiguration_id field in the database.
	FieldMisconfigurationID = "misconfiguration_id"
	// FieldMisconfigurationURLDetails holds the string denoting the misconfiguration_url_details field in the database.
	FieldMisconfigurationURLDetails = "misconfiguration_url_details"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldSeverity holds the string denoting the severity field in the database.
	FieldSeverity = "severity"
	// EdgeManifestMisconfigurations holds the string denoting the manifest_misconfigurations edge name in mutations.
	EdgeManifestMisconfigurations = "manifest_misconfigurations"
	// Table holds the table name of the misconfiguration in the database.
	Table = "misconfigurations"
	// ManifestMisconfigurationsTable is the table that holds the manifest_misconfigurations relation/edge.
	ManifestMisconfigurationsTable = "manifest_misconfigurations"
	// ManifestMisconfigurationsInverseTable is the table name for the ManifestMisconfiguration entity.
	// It exists in this package in order to avoid circular dependency with the "manifestmisconfiguration" package.
	ManifestMisconfigurationsInverseTable = "manifest_misconfigurations"
	// ManifestMisconfigurationsColumn is the table column denoting the manifest_misconfigurations relation/edge.
	ManifestMisconfigurationsColumn = "misconfiguration_id"
)

// Columns holds all SQL columns for misconfiguration fields.
var Columns = []string{
	FieldID,
	FieldMisconfigurationID,
	FieldMisconfigurationURLDetails,
	FieldTitle,
	FieldSeverity,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Severity defines the type for the "severity" enum field.
type Severity string

// Severity values.
const (
	SeverityUNKNOWN  Severity = "UNKNOWN"
	SeverityLOW      Severity = "LOW"
	SeverityMEDIUM   Severity = "MEDIUM"
	SeverityHIGH     Severity = "HIGH"
	SeverityCRITICAL Severity = "CRITICAL"
)

func (s Severity) String() string {
	return string(s)
}

// SeverityValidator is a validator for the "severity" field enum values. It is called by the builders before save.
func SeverityValidator(s Severity) error {
	switch s {
	case SeverityUNKNOWN, SeverityLOW, SeverityMEDIUM, SeverityHIGH, SeverityCRITICAL:
		return nil
	default:
		return fmt.Errorf("misconfiguration: invalid enum value for severity field: %q", s)
	}
}

// OrderOption defines the ordering options for the Misconfiguration queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByMisconfigurationID orders the results by the misconfiguration_id field.
func ByMisconfigurationID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMisconfigurationID, opts...).ToFunc()
}

// ByMisconfigurationURLDetails orders the results by the misconfiguration_url_details field.
func ByMisconfigurationURLDetails(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMisconfigurationURLDetails, opts...).ToFunc()
}

// ByTitle orders the results by the title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// BySeverity orders the results by the severity field.
func BySeverity(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSeverity, opts...).ToFunc()
}

// ByManifestMisconfigurationsCount orders the results by manifest_misconfigurations count.
func ByManifestMisconfigurationsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newManifestMisconfigurationsStep(), opts...)
	}
}

// ByManifestMisconfigurations orders the results by manifest_misconfigurations terms.
func ByManifestMisconfigurations(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newManifestMisconfigurationsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newManifestMisconfigurationsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ManifestMisconfigurationsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ManifestMisconfigurationsTable, ManifestMisconfigurationsColumn),
	)
}
