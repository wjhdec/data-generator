package schema

type Schema struct {
	FieldIndex       int
	Field, FieldType string
	Length           int
	NotNull          bool
}

type SchemaReader interface {
	// Read 读取表的字段信息
	Read(table string) []*Schema
}
