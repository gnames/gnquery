package search_test

import (
	"testing"

	"github.com/gnames/gnquery/ent/search"
	"github.com/stretchr/testify/assert"
)

func TestToQuery(t *testing.T) {
	tests := []struct {
		msg string
		inp search.Input
		q   string
	}{
		{
			"with query",
			search.Input{Query: "some query", DataSources: []int{5, 4}, NameString: "Bubo bubo"},
			"ds:5,4 n:Bubo bubo",
		},
		{
			"with query",
			search.Input{Query: "some query", DataSources: []int{5, 4}, NameString: "Bubo bubo"},
			"ds:5,4 n:Bubo bubo",
		},
		{
			"with all res",
			search.Input{DataSources: []int{5,4}, NameString: "Bubo bubo", WithAllMatches: true},
			"ds:5,4 n:Bubo bubo all:true",
		},
		{
			"with name string",
			search.Input{DataSources: []int{5}, NameString: "Bubo bubo", Author: "L."},
			"ds:5 n:Bubo bubo",
		},
		{
			"data1",
			search.Input{DataSources: []int{5}, Genus: "Bubo", SpeciesAny: "bubo", Author: "L."},
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
