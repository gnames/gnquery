package search

import "github.com/gnames/gnlib/ent/verifier"

type Output struct {
	Meta  `json:"metadata"`
	Names []Canonical `json:"results,omitempty"`
}

type Meta struct {
	Input
}

type Canonical struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Cardinality int                    `json:"cardinality,omitempty"`
	MatchType   string                 `json:"matchType"`
	BestResult  *verifier.ResultData   `json:"bestResult,omitempty"`
	Results     []*verifier.ResultData `json:"preferredResults,omitempty"`
}
