package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Consumer holds the schema definition for the Consumer entity.
type Consumer struct {
	ent.Schema
}

// Fields of the Consumer.
func (Consumer) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.Bool("submission_status").Default(false),
		field.String("username").Optional(),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Optional(),
		field.Int64("uploadfile_id").Optional(),
	}
}

// Edges of the Consumer.
func (Consumer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("uploadedfile", UploadedFile.Type).Ref("consumer").Unique().Field("uploadfile_id"),
	}
}
