package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

var TimeStampWithTZ = map[string]string{dialect.Postgres: "TIMESTAMP(0) WITH TIME ZONE"}

// SERP holds the schema definition for the SERPs entity.
type SERP struct {
	ent.Schema
}

// Annotations of the SERP.
func (SERP) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "serps"},
	}
}

// Fields of the SERP.
func (SERP) Fields() []ent.Field {
	return []ent.Field{
		field.String("url").NotEmpty(),
		field.String("title").NotEmpty(),
		field.String("description"),
		field.JSON("phones", []string{}).Optional(),
		field.JSON("emails", []string{}).Optional(),
		field.JSON("key_words", []string{}),
		field.Bool("is_read").Default(false),
		field.Int("sq_id").Optional(),
		field.Time("created_at").SchemaType(TimeStampWithTZ).Default(time.Now),
	}
}

// Edges of the SERP.
func (SERP) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("search_query", SearchQuery.Type).
			Ref("serps").
			Unique().
			Field("sq_id"),
	}
}

// Indexes of the SERP.
func (SERP) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("url", "phones", "emails", "sq_id").
			Unique(),
	}
}
