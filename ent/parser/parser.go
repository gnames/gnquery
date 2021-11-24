package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gnames/gnquery/ent/search"
	"github.com/gnames/gnquery/ent/tag"
)

type parser struct {
	elements map[tag.Tag]string
	warnings []string
	tail     string
	*engine
}

// New creates a new Parser object for converting a query string into
// search.Input object.
func New() Parser {
	res := parser{engine: &engine{}, elements: make(map[tag.Tag]string)}
	res.Init()
	return &res
}

func (p *parser) ParseQuery(q string) search.Input {
	var err error
	p.Buffer = q
	p.resetFull()
	res := search.Input{}
	if err = p.Parse(); err != nil {
		errMsg := fmt.Sprintf("Could not finish query parsing: %s", err)
		res.Warnings = append(res.Warnings, errMsg)
		return res
	}
	p.walkAST()
	return p.query()
}

// Debug takes a string, parses it, and prints its AST.
func (p *parser) Debug(q string) {
	p.Buffer = q
	p.resetFull()
	err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}
	p.PrettyPrintSyntaxTree(q)
}

func (p *parser) resetFull() {
	p.elements = make(map[tag.Tag]string)
	p.tail = ""
	p.reset()
}

func (p *parser) query() search.Input {
	var warn string
	res := search.Input{
		Query:        p.Buffer,
		Warnings:     p.warnings,
		NameString:   p.elements[tag.NameString],
		Uninomial:    p.elements[tag.Uninomial],
		Genus:        p.elements[tag.Genus],
		ParentTaxon:  p.elements[tag.ParentTaxon],
		Species:      p.elements[tag.Species],
		SpeciesAny:   p.elements[tag.SpeciesAny],
		SpeciesInfra: p.elements[tag.SpeciesInfra],
		Author:       p.elements[tag.Author],
		Tail:         p.tail,
	}

	if res.Tail != "" {
		res.Warnings = append(res.Warnings, "Unparsed tail")
	}

	if res.Uninomial != "" && res.Genus != "" {
		res.Warnings = append(res.Warnings, "Genus and uninomial tags are incompatible")
	}

	if res.Uninomial != "" && (res.Species+res.SpeciesAny+res.SpeciesInfra != "") {
		res.Warnings = append(res.Warnings, "Species and uninomial tags are incompatible")
	}

	if res.NameString != "" && (res.Uninomial+res.Genus+res.Species+res.SpeciesAny+res.SpeciesInfra != "") {
		res.Warnings = append(res.Warnings, "If name-string is given, uninomial, genus, species tags are ignored")
	}

	dsStr := p.elements[tag.DataSourceID]
	if ds, err := strconv.Atoi(p.elements[tag.DataSourceID]); err == nil {
		res.DataSourceID = ds
	} else if dsStr != "" {
		warn = fmt.Sprintf("Cannot convert dataSourceId from '%s'", dsStr)
		res.Warnings = append(res.Warnings, warn)
	}

	yrStr := p.elements[tag.Year]
	if yr, err := strconv.Atoi(yrStr); err == nil {
		res.Year = yr
	} else if yrStr != "" {
		warn = fmt.Sprintf("Cannot convert Year from '%s'", yrStr)
		res.Warnings = append(res.Warnings, warn)
	}
	return res.ProcessName()
}
