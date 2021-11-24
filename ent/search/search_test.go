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

func TestToQuery(t *testing.T) {
	tests := []struct {
		msg string
		inp search.Input
		q   string
	}{
		{
			"with query",
			search.Input{Query: "some query", DataSourceID: 5, NameString: "Bubo bubo"},
			"some query",
		},
		{
			"with name string",
			search.Input{DataSourceID: 5, NameString: "Bubo bubo", Author: "L."},
			"ds:5 n:Bubo bubo",
		},
		{
			"data1",
			search.Input{DataSourceID: 5, Genus: "Bubo", SpeciesAny: "bubo", Author: "L."},
			"ds:5 g:Bubo sp+:bubo au:L.",
		},
		{
			"data2",
			search.Input{ParentTaxon: "Aves", Genus: "Bubo", SpeciesInfra: "bubo", Author: "L.", Year: 1888},
			"tx:Aves g:Bubo isp:bubo au:L. y:1888",
		},
	}

	for _, v := range tests {
		res := v.inp.ToQuery()
		assert.Equal(t, res, v.q)
	}
}
