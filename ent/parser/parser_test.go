package parser_test

import (
	"testing"

	"github.com/gnames/gnquery/ent/parser"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	q := "n:Bubo bubo sp:bubo ds:1"
	p := &parser.Engine{Buffer: q}
	p.Init()
	err := p.Parse()
	assert.Nil(t, err)
	p.Debug(q)
}
