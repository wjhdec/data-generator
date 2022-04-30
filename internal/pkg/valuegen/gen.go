package valuegen

import (
	"math/rand"
	"strconv"
	"time"
)

type ValueGen interface {
	// Value 生成创建的sql字符串
	Value() string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type StringGen struct {
	// Length 字符串长度
	Length int `json:"length,omitempty"`
	// Chars 限制生成随机字符串的可用字符
	Chars []rune `json:"chars,omitempty"`
	// In 随机选择 in 里面的内容作为返回值
	In []string `json:"in,omitempty"`
}

func (g *StringGen) Value() string {
	if len(g.In) > 0 {
		return g.In[rand.Intn(len(g.In))]
	}
	length := 50
	if g.Length > 0 {
		length = g.Length
	}
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	if len(g.Chars) > 0 {
		chars = g.Chars
	}
	charsSize := len(chars)
	l := make([]rune, length)
	for i := range l {
		l[i] = chars[rand.Intn(charsSize)]
	}
	return string(l)
}

type Int32Gen struct {
	In  []int32 `json:"in,omitempty"`
	Max int32   `json:"max,omitempty"`
	Min int32   `json:"min,omitempty"`
}

func (g *Int32Gen) Value() string {
	var value int32
	if len(g.In) > 0 {
		value = g.In[rand.Intn(len(g.In))]
	} else if g.Max > g.Min {
		value = rand.Int31n(g.Max) + g.Min
	} else {
		value = rand.Int31()
	}
	return strconv.Itoa(int(value))
}

type Int64Gen struct {
	In  []int64 `json:"in,omitempty"`
	Max int64   `json:"max,omitempty"`
	Min int64   `json:"min,omitempty"`
}

func (g *Int64Gen) Value() string {
	var value int64
	if len(g.In) > 0 {
		value = g.In[rand.Intn(len(g.In))]
	} else if g.Max > g.Min {
		value = rand.Int63n(g.Max) + g.Min
	} else {
		value = rand.Int63()
	}
	return strconv.FormatInt(value, 10)
}

type DoubleGen struct {
	In  []float64 `json:"in,omitempty"`
	Max float64   `json:"max,omitempty"`
	Min float64   `json:"min,omitempty"`
}

func (g *DoubleGen) Value() string {
	var v float64
	if len(g.In) > 0 {
		v = g.In[rand.Intn(len(g.In))]
	} else if g.Max > g.Min {
		v = rand.Float64()*(g.Max-g.Min) + g.Min
	} else {
		v = rand.Float64()
	}
	return strconv.FormatFloat(v, 'f', 5, 64)
}

type TimeGen struct {
	In     []time.Time
	Max    time.Time
	Min    time.Time
	Format string
}

func (g *TimeGen) Value() string {
	format := "2006-01-02 15:04:05"
	if g.Format != "" {
		format = g.Format
	}
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.Local)
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.Local)
	if g.Max.After(g.Min) {
		min = g.Min
		max = g.Max
	}
	delta := max.Unix() - min.Unix()
	sec := rand.Int63n(delta) + min.Unix()
	return time.Unix(sec, 0).Format(format)
}
