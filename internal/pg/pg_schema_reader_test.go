package pg

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetSchema(t *testing.T) {
	a := assert.New(t)
	viper.AddConfigPath("../../configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}
	pool, err := pgxpool.Connect(context.Background(), viper.GetString("database.dsn"))
	if err != nil {
		log.Panic(err)
	}
	schemaReader := &SchemaReader{pool: pool}
	schemas := schemaReader.Read(viper.GetString("table.name"))
	for _, v := range schemas {
		t.Log(v)
	}
	a.NotEmpty(schemas, "输出不应该为空")
}
