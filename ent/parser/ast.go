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
	for k, v := range p.nameElements {
		p.setElement(k, v)
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
	tagNode := n.up
	valNode := tagNode.next
	val := p.nodeString(valNode)
	switch tagNode.pegRule {
	case ruleNameString:
		p.setElement(tag.NameString, val)
		p.newNameString(tagNode)
	case ruleAuthor:
		p.setElement(tag.Author, val)
	case ruleYear:
		p.setElement(tag.Year, val)
	case ruleDataSources:
		p.setElement(tag.DataSourceIDs, val)
	case ruleGenus:
		p.setElement(tag.Genus, val)
	case ruleParentTaxon:
		p.setElement(tag.ParentTaxon, val)
	case ruleSpecies:
		p.setElement(tag.Species, val)
	case ruleSpeciesAny:
		p.setElement(tag.SpeciesAny, val)
	case ruleSpeciesInfra:
		p.setElement(tag.SpeciesInfra, val)
	case ruleAllMatches:
		p.setElement(tag.AllMatches, val)
	case ruleAllBestResults:
		p.setElement(tag.AllBestResults, val)
	}
}

func (p *parser) newNameString(n *node32) {
	n = n.next.up
	for n != nil {
		val := p.nodeString(n)
		switch n.pegRule {
		case ruleAuVal:
			p.nameElements[tag.Author] = val
		case ruleYearVal:
			p.nameElements[tag.Year] = val
		case ruleGenusVal:
			p.nameElements[tag.Genus] = val
		case ruleSpeciesVal:
			p.nameElements[tag.Species] = val
		case ruleSpeciesInfraVal:
			p.nameElements[tag.SpeciesInfra] = val
		}
		n = n.next
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
