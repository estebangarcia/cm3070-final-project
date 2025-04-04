// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/blobchunk"
)

// BlobChunk is the model entity for the BlobChunk schema.
type BlobChunk struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// UploadID holds the value of the "upload_id" field.
	UploadID string `json:"upload_id,omitempty"`
	// SessionID holds the value of the "session_id" field.
	SessionID string `json:"session_id,omitempty"`
	// RangeFrom holds the value of the "range_from" field.
	RangeFrom uint64 `json:"range_from,omitempty"`
	// RangeTo holds the value of the "range_to" field.
	RangeTo uint64 `json:"range_to,omitempty"`
	// PartNumber holds the value of the "part_number" field.
	PartNumber   uint64 `json:"part_number,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*BlobChunk) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case blobchunk.FieldID, blobchunk.FieldRangeFrom, blobchunk.FieldRangeTo, blobchunk.FieldPartNumber:
			values[i] = new(sql.NullInt64)
		case blobchunk.FieldUploadID, blobchunk.FieldSessionID:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the BlobChunk fields.
func (bc *BlobChunk) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case blobchunk.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			bc.ID = int(value.Int64)
		case blobchunk.FieldUploadID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field upload_id", values[i])
			} else if value.Valid {
				bc.UploadID = value.String
			}
		case blobchunk.FieldSessionID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field session_id", values[i])
			} else if value.Valid {
				bc.SessionID = value.String
			}
		case blobchunk.FieldRangeFrom:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field range_from", values[i])
			} else if value.Valid {
				bc.RangeFrom = uint64(value.Int64)
			}
		case blobchunk.FieldRangeTo:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field range_to", values[i])
			} else if value.Valid {
				bc.RangeTo = uint64(value.Int64)
			}
		case blobchunk.FieldPartNumber:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field part_number", values[i])
			} else if value.Valid {
				bc.PartNumber = uint64(value.Int64)
			}
		default:
			bc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the BlobChunk.
// This includes values selected through modifiers, order, etc.
func (bc *BlobChunk) Value(name string) (ent.Value, error) {
	return bc.selectValues.Get(name)
}

// Update returns a builder for updating this BlobChunk.
// Note that you need to call BlobChunk.Unwrap() before calling this method if this BlobChunk
// was returned from a transaction, and the transaction was committed or rolled back.
func (bc *BlobChunk) Update() *BlobChunkUpdateOne {
	return NewBlobChunkClient(bc.config).UpdateOne(bc)
}

// Unwrap unwraps the BlobChunk entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (bc *BlobChunk) Unwrap() *BlobChunk {
	_tx, ok := bc.config.driver.(*txDriver)
	if !ok {
		panic("ent: BlobChunk is not a transactional entity")
	}
	bc.config.driver = _tx.drv
	return bc
}

// String implements the fmt.Stringer.
func (bc *BlobChunk) String() string {
	var builder strings.Builder
	builder.WriteString("BlobChunk(")
	builder.WriteString(fmt.Sprintf("id=%v, ", bc.ID))
	builder.WriteString("upload_id=")
	builder.WriteString(bc.UploadID)
	builder.WriteString(", ")
	builder.WriteString("session_id=")
	builder.WriteString(bc.SessionID)
	builder.WriteString(", ")
	builder.WriteString("range_from=")
	builder.WriteString(fmt.Sprintf("%v", bc.RangeFrom))
	builder.WriteString(", ")
	builder.WriteString("range_to=")
	builder.WriteString(fmt.Sprintf("%v", bc.RangeTo))
	builder.WriteString(", ")
	builder.WriteString("part_number=")
	builder.WriteString(fmt.Sprintf("%v", bc.PartNumber))
	builder.WriteByte(')')
	return builder.String()
}

// BlobChunks is a parsable slice of BlobChunk.
type BlobChunks []*BlobChunk
