package parser

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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

	if res.NameString != "" && (res.Genus+res.Species+res.SpeciesAny+res.SpeciesInfra != "") {
		res.Warnings = append(res.Warnings, "If name-string is given, genus, species tags are ignored")
	}

	dsStr := p.elements[tag.DataSourceID]
	if ds, err := strconv.Atoi(p.elements[tag.DataSourceID]); err == nil {
		res.DataSourceID = ds
	} else if dsStr != "" {
		warn = fmt.Sprintf("Cannot convert dataSourceId from '%s'", dsStr)
		res.Warnings = append(res.Warnings, warn)
	}

	yrStr := p.elements[tag.Year]
	yr, yrRange := processYear(yrStr)

	if yrStr != "" && yr == 0 && yrRange == nil {
		warn = fmt.Sprintf("Cannot convert Year from '%s'", yrStr)
		res.Warnings = append(res.Warnings, warn)
	}

	if yrRange != nil && yrRange.YearStart+yrRange.YearEnd == 0 {
		warn = fmt.Sprintf("Cannot convert YearRange from '%s'", yrStr)
		res.Warnings = append(res.Warnings, warn)
	}

	res.Year = yr
	res.YearRange = yrRange

	return res.ProcessName()
}

func processYear(val string) (int, *search.YearRange) {
	var yr, yrStart, yrEnd int

	yrs := strings.Split(val, "-")
	if len(yrs) == 2 {
		yrStart, _ = strconv.Atoi(yrs[0])
		yrEnd, _ = strconv.Atoi(yrs[1])
		yRange := &search.YearRange{YearStart: yrStart, YearEnd: yrEnd}
		return yr, yRange
	}

	yr, _ = strconv.Atoi(val)
	return yr, nil
}
