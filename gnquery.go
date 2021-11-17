package gnquery

import "github.com/gnames/gnquery/ent/query"

type gnquery struct{}

func New() GNquery {
	return gnquery{}
}

func (gnq gnquery) Parse(q string) query.Query {
	var res query.Query
	return res
}
