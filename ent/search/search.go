package search

import "github.com/gnames/gnlib/ent/verifier"

type Output struct {
	Meta    `json:"metadata"`
	Results []Canonical `json:"results,omitempty"`
}

type Meta struct {
	Input
}

type Canonical struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	MatchType        string                 `json:"matchType"`
	BestResult       *verifier.ResultData   `json:"bestResult,omitempty"`
	PreferredResults []*verifier.ResultData `json:"preferredResults,omitempty"`
}
