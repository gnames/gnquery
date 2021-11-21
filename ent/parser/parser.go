package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gnames/gnquery/ent/query"
)

type parser struct {
	elements map[Tag]string
	warnings []string
	tail     string
	*engine
}

func New() Parser {
	res := parser{engine: &engine{}, elements: make(map[Tag]string)}
	res.Init()
	return &res
}

func (p *parser) ParseQuery(q string) query.Query {
	var err error
	p.Buffer = q
	p.resetFull()
	res := query.Query{}
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
	p.elements = make(map[Tag]string)
	p.tail = ""
	p.reset()
}

func (p *parser) query() query.Query {
	var warn string
	res := query.Query{
		Input:        p.Buffer,
		Warnings:     p.warnings,
		NameString:   p.elements[NameString],
		Uninomial:    p.elements[Uninomial],
		Genus:        p.elements[Genus],
		ParentTaxon:  p.elements[ParentTaxon],
		Species:      p.elements[Species],
		SpeciesAny:   p.elements[SpeciesAny],
		SpeciesInfra: p.elements[SpeciesInfra],
		Author:       p.elements[Author],
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

	dsStr := p.elements[DataSourceID]
	if ds, err := strconv.Atoi(p.elements[DataSourceID]); err == nil {
		res.DataSourceID = ds
	} else if dsStr != "" {
		warn = fmt.Sprintf("Cannot convert dataSourceId from '%s'", dsStr)
		res.Warnings = append(res.Warnings, warn)
	}

	yrStr := p.elements[Year]
	if yr, err := strconv.Atoi(yrStr); err == nil {
		res.Year = yr
	} else if yrStr != "" {
		warn = fmt.Sprintf("Cannot convert Year from '%s'", yrStr)
		res.Warnings = append(res.Warnings, warn)
	}
	return res.ProcessName()
}
