package parser_test

import (
	"log"
	"testing"

	"github.com/gnames/gnquery/ent/parser"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	q := "n:Bubo bubo gen:Bubo sp+:bubo ds:1 y:-"
	parser.New().Debug(q)
}

func TestParseQuery(t *testing.T) {
	s := "n:Bubo bubo sp:caboom ds:1"
	p := parser.New()
	q := p.ParseQuery(s)
	assert.True(t, len(q.Warnings) > 0)
	assert.Equal(t, q.Query, "n:Bubo bubo sp:caboom ds:1")
	assert.Equal(t, q.NameString, "Bubo bubo")
	assert.Equal(t, q.Species, "bubo")
	assert.Equal(t, q.DataSourceID, 1)
}

func TestYearRange(t *testing.T) {
	tests := []struct {
		msg, q string
		warns  int
		yst    int
		yend   int
	}{
		{"full range", "g:B. sp:b. y:1888-2000", 0, 1888, 2000},
		{"range", "g:B. sp:b. y:1888-2000", 0, 1888, 2000},
		{"greater", "g:B. sp:b. y:1888-", 0, 1888, 0},
		{"less", "g:B. sp:b. y:-2000", 0, 0, 2000},
	}

	p := parser.New()

	for _, v := range tests {
		res := p.ParseQuery(v.q)
		assert.True(t, len(res.Warnings) == v.warns)
		assert.Equal(t, res.Year, 0)
		assert.NotNil(t, res.YearRange)
		assert.Equal(t, res.YearStart, v.yst)
		assert.Equal(t, res.YearEnd, v.yend)
	}

	q := "g:B. sp:b. y:-"
	res := p.ParseQuery(q)
	assert.True(t, len(res.Warnings) > 0)
	log.Printf("P: %#v", res)
	assert.Nil(t, res.YearRange)
}
