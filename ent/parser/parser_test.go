package parser_test

import (
	"testing"

	"github.com/gnames/gnquery/ent/parser"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	q := "n:Bubo bubo au:Linn. sp:caboom all:t ds:1"
	parser.New().Debug(q)
}

func TestParseQuery(t *testing.T) {
	s := "n:Bubo bubo au:Linn. sp:caboom all:t ds:1,2"
	p := parser.New()
	q := p.ParseQuery(s)
	assert.True(t, len(q.Warnings) > 0, "warn")
	assert.Equal(t, q.Query, "n:Bubo bubo au:Linn. sp:caboom all:t ds:1", "query")
	assert.Equal(t, q.NameString, "Bubo bubo", "n")
	assert.Equal(t, q.Species, "bubo", "sp")
	assert.Equal(t, q.DataSourceIDs, []int{1, 2}, "ds")
	assert.Equal(t, q.WithAllResults, true, "all")
	assert.Equal(t, q.Author, "Linn.", "au")
}

func TestQueries(t *testing.T) {
	tests := []struct {
		msg, q string
		warns  int
		yst    int
		yend   int
		ds     []int
		all    bool
	}{
		{"full range", "g:B. sp:b. y:1888-2000", 0, 1888, 2000, []int{}, false},
		{"range", "g:B. sp:b. y:1888-2000", 0, 1888, 2000, []int{}, false},
		{"greater", "g:B. sp:b. y:1888-", 0, 1888, 0, []int{}, false},
		{"less", "g:B. sp:b. y:-2000", 0, 0, 2000, []int{}, false},
		{"all", "g:B. sp:b. y:-2000 all:t", 0, 0, 2000, []int{}, true},
		{"all2", "g:B. sp:b. y:-2000 all:true", 0, 0, 2000, []int{}, true},
		{"best only", "g:B. sp:b. y:-2000 all:f", 0, 0, 2000, []int{}, false},
		{"best only", "g:B. sp:b. y:-2000 all:false", 0, 0, 2000, []int{}, false},
		{"mult ds", "g:B. sp:b. ds:1,2,3 y:-2000 all:false", 0, 0, 2000, []int{1, 2, 3}, false},
		{"single ds", "g:B. sp:b. ds:1 y:-2000 all:false", 0, 0, 2000, []int{1}, false},
	}

	p := parser.New()

	for _, v := range tests {
		res := p.ParseQuery(v.q)
		assert.True(t, len(res.Warnings) == v.warns)
		assert.Equal(t, res.WithAllResults, v.all)
		assert.Equal(t, res.Year, 0)
		assert.NotNil(t, res.YearRange)
		assert.Equal(t, res.YearStart, v.yst)
		assert.Equal(t, res.YearEnd, v.yend)
	}

	q := "g:B. sp:b. y:-"
	res := p.ParseQuery(q)
	assert.True(t, len(res.Warnings) > 0)
	assert.Nil(t, res.YearRange)
}
