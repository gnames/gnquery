package search

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/gnames/gnquery/ent/tag"
)

// Input contains data needed for creating a faceted search query by
// GNames.
type Input struct {
	// Query field contains the query of the faceted search.
	//
	// Example: `g:Bubo sp:bubo a:Linn. yr:1700-1800`
	Query string `json:"query,omitempty"`

	// DataSourceID field contains data-source ID for ParentTaxon and
	// serves as a preferred data-source in the final search results.
	DataSources []int `json:"dataSources,omitempty"`

	// WithAllMatches field indicates if all found data should be provided by the
	// output, or only the "best result".
	WithAllMatches bool `json:"withAllMatches,omitempty"`

	// ParentTaxon creates a filter to return only results that have the
	// required taxon in their classification path. The field uses the first
	// ID from DataSourceIDs.
	ParentTaxon string `json:"parentTaxon,omitempty"`

	// NameString is a convenience field. It allows to provide faceted search
	// data as a freeform 'name-string' input. The name-string does not have
	// to be a 'real name', and it suppose to have only one component of each
	// kind. For example it should contain one infraspecific epithet, and
	// one author. If it has more, a warning will be generated and only one
	// of the values of a ducplicated category will be used.
	//
	// For example: `Aus bus cus A. 1888`.
	//
	// If NameString is given, all other inputs for genus, species,
	// authors, year will be ignored.
	NameString string `json:"nameString,omitempty"`

	// Genus field creates a faceted search filter by a genus name. The name
	// can be abbreviated up to one character following a period.
	//
	// For example: `H.`, `Hom.`, `Homo`.
	Genus string `json:"genus,omitempty"`

	// SpeciesAny field creates a faceted search filter for either specific or
	// infraspecific epithets. If SpeciesAny is given Species and SpeciesInfra
	// fields are ignored.
	SpeciesAny string `json:"speciesAny,omitempty"`

	// Species field creates a faceted search filter for specific epithet only.
	Species string `json:"species,omitempty"`

	// SpecieInfra field creates a faceted search filter for infraspecific
	// epithets. Such epithets can be subspecies, varieties, formas etc.
	SpeciesInfra string `json:"speciesInfra,omitempty"`

	// Author field creates a faceted search filter for Author name. This
	// field can be abbreviated up to one letter folloing period character.
	//
	// For example: `L.`, `Linn.`, `Linnaeus`.
	Author string `json:"author,omitempty"`

	// Year field creates a faceted search filter for a year.
	Year int `json:"year,omitempty"`

	// YearRange field creates a more flexible filter for year data.
	//
	// Example:
	// 1990-1995: year >= 1990 && year =< 1995
	// 1990-: year >= 1990
	// -1995: year <= 1995
	*YearRange `json:"yearRange,omitempty"`

	// Tail field keeps a non-parsed pard of a query, in case if full parsing
	// of query string did fail.
	Tail string `json:"unparsedTail,omitempty"`

	// Warnings field contains warnings generated by a parsing process.
	Warnings []string `json:"warnings,omitempty"`
}

// YearRange field creates a more flexible filter for year data.
//
// Example:
// 1990-1995: year >= 1990 && year =< 1995
// 1990-: year >= 1990
// -1995: year <= 1995
type YearRange struct {
	// YearStart is the smaller year of the range.
	// YearEnd is the larger year of the range.
	YearStart, YearEnd int
}

type qTags struct {
	tag tag.Tag
	val string
}

// ToQuery creates a query from the Input data.
func (inp Input) ToQuery() string {
	qSlice := make([]string, 0, 9)

	data1 := []qTags{
		{tag.DataSourceIDs, inp.dsToString()},
		{tag.ParentTaxon, inp.ParentTaxon},
		{tag.NameString, inp.NameString},
		{tag.AllMatches, strconv.FormatBool(inp.WithAllMatches)},
	}

	qSlice = processTags(qSlice, data1)

	if inp.NameString != "" {
		return strings.Join(qSlice, " ")
	}

	var yr, yrStart, yrEnd string
	if val := strconv.Itoa(inp.Year); val != "0" {
		yr = val
	}

	if inp.YearRange != nil {
		if inp.YearStart > 0 {
			yrStart = strconv.Itoa(inp.YearStart)
		}
		if inp.YearEnd > 0 {
			yrEnd = strconv.Itoa(inp.YearEnd)
		}
		yr = fmt.Sprintf("%s-%s", yrStart, yrEnd)
	}

	data2 := []qTags{
		{tag.Genus, inp.Genus},
		{tag.Species, inp.Species},
		{tag.SpeciesInfra, inp.SpeciesInfra},
		{tag.SpeciesAny, inp.SpeciesAny},
		{tag.Author, inp.Author},
		{tag.Year, yr},
	}

	qSlice = processTags(qSlice, data2)

	return strings.Join(qSlice, " ")
}

func processTags(qSlice []string, data []qTags) []string {
	for _, v := range data {
		if v.val != "" && v.val != "false" {
			str := v.tag.String() + ":" + v.val
			qSlice = append(qSlice, str)
		}
	}
	return qSlice
}

func (inp Input) dsToString() string {
	if len(inp.DataSources) == 0 {
		return ""
	}

	res := make([]string, len(inp.DataSources))
	for i := range inp.DataSources {
		res[i] = strconv.Itoa(inp.DataSources[i])
	}
	return strings.Join(res, ",")
}

// IsQuery is a very simple determination if a string looks like a
// faceted search query. If the string starts with low-case letters following
// a colon, the string is considered to be a query.
//
// Such a simple method is used, because neither scientific names or
// files normally have such pattern at th start of their string.
func IsQuery(s string) bool {
	s = strings.TrimSpace(s)
	idx := strings.Index(s, ":")
	if idx == -1 {
		return false
	}

	rs := []rune(s[0:idx])
	for i := range rs {
		if !unicode.IsLower(rs[i]) {
			return false
		}
	}
	rs = []rune(s[idx:])
	if len(rs) < 2 {
		return false
	}

	return unicode.IsLetter(rs[1]) || unicode.IsNumber(rs[1])
}
