package gnquery

import "github.com/gnames/gnquery/ent/search"
import "github.com/gnames/gnquery/ent/parser"

type gnquery struct{
  parser.Parser
}

// New returns GNquery object. GNquery provides conversion of a query into
// search parameters.
func New() GNquery {
	return gnquery{
    Parser: parser.New(),
  }
}

// Parse takes a query string and returns back a search.Input object that
// contains information to create faceted search filters by GNames.
func (gnq gnquery) Parse(q string) search.Input {
	return gnq.ParseQuery(q)
}
