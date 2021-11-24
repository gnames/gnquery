package parser

import (
	"fmt"
	"strings"

	"github.com/gnames/gnquery/ent/tag"
)

func (p *parser) walkAST() {
	p.newQuery(p.AST())
}

func (p *parser) newQuery(n *node32) {
	n = n.up

	for n != nil {
		switch n.pegRule {
		case ruleComponents:
			p.newComponents(n)
		case ruleTail:
			p.tail = p.tailValue(n)
		}
		n = n.next
	}
}

func (p *parser) newComponents(n *node32) {
	n = n.up

	for n != nil {
		switch n.pegRule {
		case ruleElement:
			p.newElement(n)
		}
		n = n.next
	}
}

func (p *parser) setElement(tag tag.Tag, val string) {
	if _, ok := p.elements[tag]; ok {
		warn := fmt.Sprintf("Tag '%s' appears more than once", tag.String())
		p.warnings = append(p.warnings, warn)
	}
	p.elements[tag] = val
}

func (p *parser) newElement(n *node32) {
	n = n.up
	tagNode := n.up
	valNode := n.next
	val := p.nodeString(valNode)

	switch tagNode.pegRule {
	case ruleAuthor:
		p.setElement(tag.Author, val)
	case ruleYear:
		p.setElement(tag.Year, val)
	case ruleDataSource:
		p.setElement(tag.DataSourceID, val)
	case ruleGenus:
		p.setElement(tag.Genus, val)
	case ruleNameString:
		p.setElement(tag.NameString, val)
	case ruleParentTaxon:
		p.setElement(tag.ParentTaxon, val)
	case ruleSpecies:
		p.setElement(tag.Species, val)
	case ruleSpeciesAny:
		p.setElement(tag.SpeciesAny, val)
	case ruleSpeciesInfra:
		p.setElement(tag.SpeciesInfra, val)
	case ruleUninomial:
		p.setElement(tag.Uninomial, val)
	}
}

func (p *parser) nodeString(n *node32) string {
	t := n.token32
	return string(p.buffer[t.begin:t.end])
}

func (p *parser) tailValue(n *node32) string {
	res := p.nodeString(n)
	return strings.TrimRight(res, " ")
}
