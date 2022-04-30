package pg

import (
	"context"
	"datagen/internal/pkg/schema"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SchemaReader struct {
	pool *pgxpool.Pool
}

func (p *SchemaReader) Read(table string) []*schema.Schema {
	var schemas []*schema.Schema
	schemaSql := `
	select
		a.attnum as field_index,
		a.attname as field,
		t.typname as field_type,
		a.attlen as length,
		a.attnotnull as not_null
	from
		pg_class c,
		pg_attribute a,
		pg_type t
	where
		c.relname = '%s'
		and a.attnum > 0
		and a.attrelid = c.oid
		and a.atttypid = t.oid
	order by
		a.attnum
	`
	rows, err := p.pool.Query(context.Background(), fmt.Sprintf(schemaSql, table))
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		var (
			field, fieldType   string
			fieldIndex, length int
			notNull            bool
		)
		err := rows.Scan(&fieldIndex, &field, &fieldType, &length, &notNull)
		if err != nil {
			log.Panic(err)
		}
		s := &schema.Schema{FieldIndex: fieldIndex, Field: field, FieldType: fieldType, Length: length, NotNull: notNull}
		schemas = append(schemas, s)
	}
	if len(schemas) == 0 {
		log.Printf("can not find fields in %s", table)
	}
	return schemas
}
