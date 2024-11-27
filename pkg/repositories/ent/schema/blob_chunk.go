package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Manifest holds the schema definition for the BlobChunk entity.
type BlobChunk struct {
	ent.Schema
}

// Fields of the BlobChunk.
func (BlobChunk) Fields() []ent.Field {
	return []ent.Field{
		field.String("upload_id"),
		field.String("session_id"),
		field.Uint64("range_from"),
		field.Uint64("range_to"),
		field.Uint64("part_number"),
	}
}

func (BlobChunk) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("upload_id", "session_id", "part_number").Unique(),
		index.Fields("upload_id"),
		index.Fields("session_id"),
	}
}

// Edges of the BlobChunk.
func (BlobChunk) Edges() []ent.Edge {
	return nil
}
