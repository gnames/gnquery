package gnquery

import (
  "github.com/gnames/gnquery/ent/query"
)

type GNquery interface {
  Parse(string) query.Query
}
