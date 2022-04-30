package valuegen

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringGenSize(t *testing.T) {
	sg := &StringGen{
		Length: 10,
	}
	a := assert.New(t)
	value := sg.Value()
	t.Log(value)
	a.Len(value, 10)
}

func TestStringGenChars(t *testing.T) {
	sg := &StringGen{
		Length: 10,
		Chars:  []rune("a"),
	}
	a := assert.New(t)
	value := sg.Value()
	t.Log(value)
	a.Equal("aaaaaaaaaa", value)
}

func TestStringGenIn(t *testing.T) {
	in := []string{"a1", "a2", "b1", "b2"}
	sg := &StringGen{
		Length: 10,
		In:     in,
	}
	a := assert.New(t)
	for i := 0; i < 10; i++ {
		v := sg.Value()
		t.Log(v)
		a.Contains(in, v)
	}
}

func TestInt64Gen(t *testing.T) {
	ig := &Int64Gen{}
	v := ig.Value()
	t.Log(v)
}

func TestDoubleGen(t *testing.T) {
	dg := &DoubleGen{}
	for i := 0; i < 50; i++ {
		v := dg.Value()
		t.Log(v)
	}
}

func TestDoubleGenRange(t *testing.T) {
	min := float64(-200)
	max := float64(100)
	dg := &DoubleGen{
		Max: max,
		Min: min,
	}
	a := assert.New(t)
	for i := 0; i < 500; i++ {
		v := dg.Value()
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Panic(err)
		}
		a.Less(i, max)
		a.Greater(i, min)
	}
}

func TestTimeGen(t *testing.T) {
	from := time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC)
	to := time.Date(2022, 1, 10, 0, 0, 0, 0, time.UTC)
	tg := &TimeGen{Min: from, Max: to}
	a := assert.New(t)
	for i := 0; i < 500; i++ {
		i, err := time.Parse("2006-01-02 15:04:05", tg.Value())
		if err != nil {
			log.Panic(err)
		}
		a.Less(i, to)
		a.Greater(i, from)
	}
}
