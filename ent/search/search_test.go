package search_test

import (
	"testing"

	"github.com/gnames/gnquery/ent/search"
	"github.com/stretchr/testify/assert"
)

func TestProcessName(t *testing.T) {
	inp := search.Input{NameString: "H. sapiens", Genus: "Pomatomus", Author: "L."}

	inp = inp.ProcessName()
	assert.Equal(t, inp.Author, "")
	assert.Equal(t, inp.Genus, "H.")
	assert.Equal(t, inp.Species, "sapiens")
}
