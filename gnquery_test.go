package gnquery_test

import (
	"fmt"

	"github.com/gnames/gnquery"
)

func Example() {
	q := "ds:2 tx:Aves g:Bubo sp+:bubo y:1758"
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
	// ds:2 tx:Aves g:Bubo sp+:bubo y:1758
	// 2
	// Aves
	// Bubo
	// bubo
	// 1758
  //
}
