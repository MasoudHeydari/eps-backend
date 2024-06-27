// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// SerpsColumns holds the columns for the "serps" table.
	SerpsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "url", Type: field.TypeString},
		{Name: "title", Type: field.TypeString},
		{Name: "description", Type: field.TypeString},
		{Name: "phones", Type: field.TypeJSON, Nullable: true},
		{Name: "emails", Type: field.TypeJSON, Nullable: true},
		{Name: "key_words", Type: field.TypeJSON},
		{Name: "is_read", Type: field.TypeBool, Default: false},
		{Name: "created_at", Type: field.TypeTime, SchemaType: map[string]string{"postgres": "TIMESTAMP(0) WITH TIME ZONE"}},
		{Name: "sq_id", Type: field.TypeInt, Nullable: true},
	}
	// SerpsTable holds the schema information for the "serps" table.
	SerpsTable = &schema.Table{
		Name:       "serps",
		Columns:    SerpsColumns,
		PrimaryKey: []*schema.Column{SerpsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "serps_search_queries_serps",
				Columns:    []*schema.Column{SerpsColumns[9]},
				RefColumns: []*schema.Column{SearchQueriesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "serp_url_phones_emails",
				Unique:  true,
				Columns: []*schema.Column{SerpsColumns[1], SerpsColumns[4], SerpsColumns[5]},
			},
		},
	}
	// SearchQueriesColumns holds the columns for the "search_queries" table.
	SearchQueriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "query", Type: field.TypeString},
		{Name: "location", Type: field.TypeString},
		{Name: "language", Type: field.TypeString},
		{Name: "is_canceled", Type: field.TypeBool, Default: false},
		{Name: "created_at", Type: field.TypeTime, SchemaType: map[string]string{"postgres": "TIMESTAMP(0) WITH TIME ZONE"}},
	}
	// SearchQueriesTable holds the schema information for the "search_queries" table.
	SearchQueriesTable = &schema.Table{
		Name:       "search_queries",
		Columns:    SearchQueriesColumns,
		PrimaryKey: []*schema.Column{SearchQueriesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "searchquery_query_location_language",
				Unique:  true,
				Columns: []*schema.Column{SearchQueriesColumns[1], SearchQueriesColumns[2], SearchQueriesColumns[3]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		SerpsTable,
		SearchQueriesTable,
	}
)

func init() {
	SerpsTable.ForeignKeys[0].RefTable = SearchQueriesTable
	SerpsTable.Annotation = &entsql.Annotation{
		Table: "serps",
	}
	SearchQueriesTable.Annotation = &entsql.Annotation{
		Table: "search_queries",
	}
}
