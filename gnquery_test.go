package gnquery_test

import (
	"fmt"
	"testing"

	"github.com/gnames/gnquery"
	"github.com/gnames/gnquery/ent/search"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	inp := search.Input{
		Genus:   "P.",
		Species: "saltator",
	}
	gnq := gnquery.New()
	res := gnq.Process(inp)
	val := search.Input{
		Query:   "g:P. sp:saltator",
		Genus:   "P.",
		Species: "saltator",
	}
	assert.Equal(t, res, val)

	inp = search.Input{
		Query:   "n:Puma concolor au:Linn. all:t",
		Genus:   "P.",
		Species: "saltator",
	}
	res = gnq.Process(inp)
	val = search.Input{
		Query:          "n:Puma concolor au:Linn. all:t",
		NameString:     "Puma concolor",
		Genus:          "Puma",
		Species:        "concolor",
		Author:         "Linn.",
		WithAllResults: true,
	}
	assert.Equal(t, res, val)
}

func Example() {
	q := "ds:2 tx:Aves g:Bubo asp:bubo y:1758"
	gnq := gnquery.New()
	res := gnq.Parse(q)
	fmt.Println(res.Query)
	fmt.Println(res.DataSourceID)
	fmt.Println(res.ParentTaxon)
	fmt.Println(res.Genus)
	fmt.Println(res.SpeciesAny)
	fmt.Println(res.Year)
	fmt.Println(res.Tail)
	// Output:
	// ds:2 tx:Aves g:Bubo asp:bubo y:1758
	// 2
	// Aves
	// Bubo
	// bubo
	// 1758
	//
}
