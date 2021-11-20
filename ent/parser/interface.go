package parser

import "github.com/gnames/gnquery/ent/query"

type Parser interface {
	ParseQuery(string) (query.Query, error)
	Debug(string)
}
