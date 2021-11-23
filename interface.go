// gnquery package is a parser of a search query for GNames API. The query is
// used for faceted search that can be filtered  by abbreviated name, by
// species epithet, author, year, a parent clade.
package gnquery

import (
	"github.com/gnames/gnquery/ent/search"
)

// GNquery is an object for parsing search queries by GNames.
type GNquery interface {
  // Parse takes a query and returns back an Input object for performing
  // a faceted search in GNames.
	Parse(q string) search.Input
}
