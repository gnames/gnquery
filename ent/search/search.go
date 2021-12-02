package search

import "github.com/gnames/gnlib/ent/verifier"

type Output struct {
	Meta  `json:"metadata"`
	Names []Canonical `json:"names,omitempty"`
}

type Meta struct {
	Input    Input  `json:"input"`
	NamesNum int    `json:"namesNumber"`
	Error    string `json:"error,omitempty"`
}

type Canonical struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Cardinality int                    `json:"cardinality,omitempty"`
	MatchType   string                 `json:"matchType"`
	BestResult  *verifier.ResultData   `json:"bestResult,omitempty"`
	Results     []*verifier.ResultData `json:"results,omitempty"`
}
