package parser_test

import (
	"testing"

	"github.com/gnames/gnquery/ent/parser"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	q := "n:Bubo bubo sp:bubo ds:1"
	parser.New().Debug(q)
}

func TestParseQuery(t *testing.T) {
	s := "n:Bubo bubo sp:bubo ds:1"
	p := parser.New()
	q, err := p.ParseQuery(s)
	assert.Nil(t, err)
	assert.Equal(t, q.Species, "bubo")
	assert.Equal(t, q.DataSource, 1)
}
