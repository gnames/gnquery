package gnquery

import (
	"github.com/gnames/gnquery/ent/parser"
	"github.com/gnames/gnquery/ent/search"
)

type gnquery struct {
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

func (gnq gnquery) Process(inp search.Input) search.Input {
	if inp.Query == "" {
		inp.Query = inp.ToQuery()
		return inp
	}

	inp2 := gnq.ParseQuery(inp.Query)
	inp.WithAllResults = inp2.WithAllResults
	inp.Warnings = inp2.Warnings

	if len(inp2.DataSourceIDs) > 0 {
		inp.DataSourceIDs = inp2.DataSourceIDs
	}

	if inp2.ParentTaxon != "" {
		inp.ParentTaxon = inp2.ParentTaxon
	}

	if inp2.NameString != "" {
		inp.NameString = inp2.NameString
	}

	if inp2.Genus != "" {
		inp.Genus = inp2.Genus
	}

	if inp2.SpeciesAny != "" {
		inp.SpeciesAny = inp2.SpeciesAny
	}

	if inp2.Species != "" {
		inp.Species = inp2.Species
	}

	if inp2.SpeciesInfra != "" {
		inp.SpeciesInfra = inp2.SpeciesInfra
	}

	if inp2.Author != "" {
		inp.Author = inp2.Author
	}

	if inp2.Year > 0 {
		inp.Year = inp2.Year
	}

	if inp2.YearRange != nil {
		inp.YearRange = inp2.YearRange
	}

	if inp2.Tail != "" {
		inp.Tail = inp2.Tail
	}

	return inp
}
