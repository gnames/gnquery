package parser

import (
	"log"
	"strconv"

	"github.com/gnames/gnquery/ent/query"
)

type parser struct {
	elements map[Tag]string
	tail     string
	*engine
}

func New() Parser {
	res := parser{engine: &engine{}, elements: make(map[Tag]string)}
	res.Init()
	return &res
}

func (p *parser) ParseQuery(q string) (query.Query, error) {
	var err error
	p.Buffer = q
	p.resetFull()
	res := query.Query{}
	if err = p.Parse(); err != nil {
		return res, err
	}
	p.walkAST()
	return p.query(), nil
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
	res := query.Query{
		Uninomial:    p.elements[Uninomial],
		Genus:        p.elements[Genus],
		ParentTaxon:  p.elements[ParentTaxon],
		Species:      p.elements[Species],
		Infraspecies: p.elements[SpeciesInfra],
		Author:       p.elements[Author],
		Tail:         p.tail,
	}
	if ds, err := strconv.Atoi(p.elements[DataSource]); err == nil {
		res.DataSource = ds
	}

	if yr, err := strconv.Atoi(p.elements[Year]); err == nil {
		res.Year = yr
	}
	return res
}
