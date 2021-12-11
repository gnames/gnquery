package search_test

import (
	"testing"

	"github.com/gnames/gnquery/ent/search"
	"github.com/stretchr/testify/assert"
)

func TestProcessName(t *testing.T) {
	inp := search.Input{NameString: "H. sapiens", Genus: "Pomatomus", Author: "L."}

	inp = inp.ProcessName()
	assert.Equal(t, inp.Author, "L.")
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
			search.Input{Query: "some query", DataSourceIDs: []int{5}, NameString: "Bubo bubo"},
			"some query",
		},
		{
			"with all res",
			search.Input{DataSourceIDs: []int{5}, NameString: "Bubo bubo", WithAllResults: true},
			"ds:5 n:Bubo bubo all:true",
		},
		{
			"with name string",
			search.Input{DataSourceIDs: []int{5}, NameString: "Bubo bubo", Author: "L."},
			"ds:5 n:Bubo bubo",
		},
		{
			"data1",
			search.Input{DataSourceIDs: []int{5}, Genus: "Bubo", SpeciesAny: "bubo", Author: "L."},
			"ds:5 g:Bubo asp:bubo au:L.",
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

func TestIsQuery(t *testing.T) {
	tests := []struct {
		msg, str string
		query    bool
	}{
		{"query", "g:B.", true},
		{"empty", "", false},
		{"yr", "  y:1888 g:B.", true},
		{"file", "c:\file", false},
		{"file2", "C:\file", false},
		{"file3", "g:\file", false},
		{"weird", "g: ", false},
		{"name", "Bubo bubo", false},
	}
	for _, v := range tests {
		res := search.IsQuery(v.str)
		assert.Equal(t, v.query, res)
	}

}
