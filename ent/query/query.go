package query

import (
	"strconv"

	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/parsed"
)

type Query struct {
	Input        string   `json:"input"`
	Warnings     []string `json:"warnings,omitempty"`
	DataSourceID int      `json:"dataSourceID,omitempty"`
	ParentTaxon  string   `json:"parentTaxon,omitempty"`
	NameString   string   `json:"nameString,omitempty"`
	Uninomial    string   `json:"uninomial,omitempty"`
	Genus        string   `json:"genus,omitempty"`
	SpeciesAny   string   `json:"speciesAny,omitempty"`
	Species      string   `json:"species,omitempty"`
	SpeciesInfra string   `json:"speciesInfra,omitempty"`
	Author       string   `json:"author,omitempty"`
	Year         int      `json:"year,omitempty"`
	Tail         string   `json:"tail,omitempty"`
}

func (q Query) ProcessName() Query {
	if q.NameString == "" {
		return q
	}

	res := Query{Input: q.Input, NameString: q.NameString,
		Tail: q.Tail, Warnings: q.Warnings,
		DataSourceID: q.DataSourceID, ParentTaxon: q.ParentTaxon,
	}

	cfg := gnparser.NewConfig(gnparser.OptWithDetails(true))
	p := gnparser.New(cfg)
	pRes := p.ParseName(q.NameString)

	if !pRes.Parsed {
		res.Warnings = append(res.Warnings, "Cannot parse '%s'", q.NameString)
		return res
	}

	for _, v := range pRes.Words {
		val := v.Normalized
		switch v.Type {
		case parsed.UninomialType:
			res.Uninomial = val
		case parsed.GenusType:
			res.Genus = val
		case parsed.SpEpithetType:
			res.Species = val
		case parsed.InfraspEpithetType:
			res.SpeciesInfra = val
		case parsed.AuthorWordType:
			res.Author = val
		case parsed.YearType:
			yr, err := strconv.Atoi(val)
			if err == nil {
				res.Year = yr
			}
		}
	}
	return res
}
