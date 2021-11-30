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
	s := "n:Bubo bubo au:Linn. sp:caboom all:t ds:1"
	p := parser.New()
	q := p.ParseQuery(s)
	assert.True(t, len(q.Warnings) > 0, "warn")
	assert.Equal(t, q.Query, "n:Bubo bubo au:Linn. sp:caboom all:t ds:1", "query")
	assert.Equal(t, q.NameString, "Bubo bubo", "n")
	assert.Equal(t, q.Species, "bubo", "sp")
	assert.Equal(t, q.DataSourceID, 1, "ds")
	assert.Equal(t, q.WithAllResults, true, "all")
	assert.Equal(t, q.Author, "Linn.", "au")
}

func TestYearRange(t *testing.T) {
	tests := []struct {
		msg, q string
		warns  int
		yst    int
		yend   int
		all    bool
	}{
		{"full range", "g:B. sp:b. y:1888-2000", 0, 1888, 2000, false},
		{"range", "g:B. sp:b. y:1888-2000", 0, 1888, 2000, false},
		{"greater", "g:B. sp:b. y:1888-", 0, 1888, 0, false},
		{"less", "g:B. sp:b. y:-2000", 0, 0, 2000, false},
		{"all", "g:B. sp:b. y:-2000 all:t", 0, 0, 2000, true},
		{"all2", "g:B. sp:b. y:-2000 all:true", 0, 0, 2000, true},
		{"best only", "g:B. sp:b. y:-2000 all:f", 0, 0, 2000, false},
		{"best only", "g:B. sp:b. y:-2000 all:false", 0, 0, 2000, false},
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
