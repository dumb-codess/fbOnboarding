package schema

import (
	"fbOnboarding/enum"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type UploadedFile struct {
	ent.Schema
}

func (UploadedFile) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("key").Unique(),
		field.String("form_status").GoType(enum.FormStatus("")),
		field.String("url").Optional(),
		field.JSON("metadata", map[string]interface{}{}).Optional(),
		field.Time("created_at").Default(time.Now()),
	}
}

func (UploadedFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("consumer", Consumer.Type).Unique(),
	}
}
