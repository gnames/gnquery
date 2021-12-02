# GNquery

This is a library for parsing search query for GNames API.
It creates an Input object that is further used for faceted search.
The search.Input object is able to create search filters genus, specific
epithet, abbreviated name, author, year, parent clade, or data-source.

## Usage

```go
import (
	"fmt"

	"github.com/gnames/gnquery"
)

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
```

## License

Released under [MIT license]

[MIT license]: https://github.com/gnames/gnquery/raw/master/LICENSE

