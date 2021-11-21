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
	q := p.ParseQuery(s)
	assert.True(t, len(q.Warnings) > 0)
	assert.Equal(t, q.Input, "n:Bubo bubo sp:bubo ds:1")
	assert.Equal(t, q.NameString, "Bubo bubo")
	assert.Equal(t, q.Species, "bubo")
	assert.Equal(t, q.DataSourceID, 1)
}
