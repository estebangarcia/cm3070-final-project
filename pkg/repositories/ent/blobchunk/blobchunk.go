// Code generated by ent, DO NOT EDIT.

package blobchunk

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the blobchunk type in the database.
	Label = "blob_chunk"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUploadID holds the string denoting the upload_id field in the database.
	FieldUploadID = "upload_id"
	// FieldSessionID holds the string denoting the session_id field in the database.
	FieldSessionID = "session_id"
	// FieldRangeFrom holds the string denoting the range_from field in the database.
	FieldRangeFrom = "range_from"
	// FieldRangeTo holds the string denoting the range_to field in the database.
	FieldRangeTo = "range_to"
	// FieldPartNumber holds the string denoting the part_number field in the database.
	FieldPartNumber = "part_number"
	// Table holds the table name of the blobchunk in the database.
	Table = "blob_chunks"
)

// Columns holds all SQL columns for blobchunk fields.
var Columns = []string{
	FieldID,
	FieldUploadID,
	FieldSessionID,
	FieldRangeFrom,
	FieldRangeTo,
	FieldPartNumber,
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

// OrderOption defines the ordering options for the BlobChunk queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByUploadID orders the results by the upload_id field.
func ByUploadID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUploadID, opts...).ToFunc()
}

// BySessionID orders the results by the session_id field.
func BySessionID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSessionID, opts...).ToFunc()
}

// ByRangeFrom orders the results by the range_from field.
func ByRangeFrom(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRangeFrom, opts...).ToFunc()
}

// ByRangeTo orders the results by the range_to field.
func ByRangeTo(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRangeTo, opts...).ToFunc()
}

// ByPartNumber orders the results by the part_number field.
func ByPartNumber(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPartNumber, opts...).ToFunc()
}
