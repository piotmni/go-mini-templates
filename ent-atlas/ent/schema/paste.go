package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Paste holds the schema definition for the Paste entity.
type Paste struct {
	ent.Schema
}

// Fields of the Paste.
func (Paste) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").
			Unique().
			NotEmpty().
			MaxLen(20),
		field.String("title").
			Optional().
			MaxLen(200),
		field.Text("content").
			NotEmpty(),
		field.String("language").
			Optional().
			Default("plaintext").
			MaxLen(50),
		field.Bool("is_public").
			Default(true),
		field.Time("expires_at").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Paste.
func (Paste) Edges() []ent.Edge {
	return nil
}

// Indexes of the Paste.
func (Paste) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug"),
		index.Fields("created_at"),
	}
}
