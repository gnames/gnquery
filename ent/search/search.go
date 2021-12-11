package search

import "github.com/gnames/gnlib/ent/verifier"

type Output struct {
	Meta  `json:"metadata"`
	Names []verifier.Name `json:"names,omitempty"`
}

type Meta struct {
	Input       Input  `json:"input"`
	NamesNumber int    `json:"namesNumber"`
	Error       string `json:"error,omitempty"`
}
