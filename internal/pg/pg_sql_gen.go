package pg

import (
	"context"
	"datagen/internal/pkg/config"
	"datagen/internal/pkg/schema"
	"datagen/internal/pkg/valuegen"
	"fmt"
	"github.com/Jeffail/tunny"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/qianlnk/pgbar"
	"github.com/spf13/viper"
	"log"
	"runtime"
	"strings"
)

type SqlGen struct {
	pool  *pgxpool.Pool
	viper *viper.Viper
}

func NewSqlGen(pool *pgxpool.Pool, viper *viper.Viper) *SqlGen {
	return &SqlGen{viper: viper, pool: pool}
}

func (g *SqlGen) CreateSql() {

	genMap := g.createGenMap()
	insertFields := make([]string, 0, len(genMap))
	insertGen := make([]valuegen.ValueGen, 0, len(genMap))
	for g := range genMap {
		insertFields = append(insertFields, g)
		insertGen = append(insertGen, genMap[g])
	}
	insertSql := g.createInsertSql(insertFields)

	allCount := g.viper.GetInt("generator.count")
	partSize := g.viper.GetInt("generator.part_size")

	partCount := allCount / partSize
	partMod := allCount % partSize
	if partMod > 0 {
		partCount += 1
	}
	pgb := pgbar.New("data generate")
	b := pgb.NewBar("finished", allCount)

	type Param struct {
		addSize int
		sql     string
		pool    *pgxpool.Pool
		bar     *pgbar.Bar
	}

	pool := tunny.NewFunc(runtime.NumCPU()*5, func(i interface{}) interface{} {
		param, ok := i.(*Param)
		if !ok {
			log.Panic("input param type error")
		}
		batch := &pgx.Batch{}
		for j := 0; j < param.addSize; j++ {
			values := make([]interface{}, 0, len(insertGen))
			for _, gen := range insertGen {
				values = append(values, gen.Value())
			}
			batch.Queue(param.sql, values...)
		}
		br := param.pool.SendBatch(context.Background(), batch)
		if _, err := br.Exec(); err != nil {
			log.Panic(err)
		}
		defer br.Close()
		b.Add(param.addSize)
		return nil
	})
	defer pool.Close()

	for i := 0; i < partCount; i++ {
		currentAddSize := partSize
		if i == partCount-1 && partMod > 0 {
			currentAddSize = partMod
		}
		param := &Param{addSize: currentAddSize, sql: insertSql, pool: g.pool}
		pool.Process(param)
	}
}

func (g *SqlGen) createInsertSql(fields []string) string {
	name := g.viper.GetString("table.name")
	var sqlFields []string
	for _, f := range fields {
		sqlFields = append(sqlFields, "\""+f+"\"")
	}
	var ss []string
	for i := 0; i < len(fields); i++ {
		ss = append(ss, fmt.Sprintf("$%d", i+1))
	}
	return fmt.Sprintf("insert into %s(%s) values (%s)", name, strings.Join(sqlFields, ","), strings.Join(ss, ","))
}

func getGen(v *viper.Viper, g valuegen.ValueGen) valuegen.ValueGen {
	if err := v.Unmarshal(g, config.DateDecoder); err != nil {
		log.Panic(err)
	}
	return g
}

func genByType(cfg *viper.Viper, s *schema.Schema) valuegen.ValueGen {
	switch s.FieldType {
	case "varchar":
		return getGen(cfg, &valuegen.StringGen{})
	case "int4":
		return getGen(cfg, &valuegen.Int32Gen{})
	case "int8":
		return getGen(cfg, &valuegen.Int64Gen{})
	case "float8", "float4":
		return getGen(cfg, &valuegen.DoubleGen{})
	case "timestamp":
		return getGen(cfg, &valuegen.TimeGen{})
	default:
		log.Panicf("暂不支持类型：%s", s.FieldType)
	}
	// 理论上不会到这里
	return nil
}

func (g *SqlGen) createGenMap() map[string]valuegen.ValueGen {
	reader := &SchemaReader{pool: g.pool}
	schemas := reader.Read(g.viper.GetString("table.name"))
	gMap := make(map[string]valuegen.ValueGen)
	for _, s := range schemas {
		// 判断自定义字段信息
		cusField := g.viper.Sub("fields.custom." + s.Field)
		if cusField != nil {
			if cusField.GetBool("skip") {
				continue
			}
			if gen := genByType(cusField, s); gen != nil {
				gMap[s.Field] = gen
				continue
			}
		}
		// 如果自定义中没有设置，则查询是否有此类型的默认配置
		if defFieldCfg := g.viper.Sub("default_cfg." + s.FieldType); defFieldCfg != nil {
			if defGen := genByType(defFieldCfg, s); defGen != nil {
				gMap[s.Field] = defGen
				continue
			}
		}
		var emptyGen valuegen.ValueGen
		switch s.FieldType {
		case "varchar":
			emptyGen = &valuegen.StringGen{}
		case "int4":
			emptyGen = &valuegen.Int32Gen{}
		case "int8":
			emptyGen = &valuegen.Int64Gen{}
		case "float8", "float4":
			emptyGen = &valuegen.DoubleGen{}
		case "timestamp":
			emptyGen = &valuegen.TimeGen{}
		default:
			log.Panicf("暂不支持类型：%s", s.FieldType)
		}
		gMap[s.Field] = emptyGen
	}
	return gMap
}
