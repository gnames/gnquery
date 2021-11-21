package parser

import (
	"fmt"
	"strings"
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

func (p *parser) setElement(tag Tag, val string) {
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

	fmt.Printf("TAG: %#v, %s\n", tagNode, p.nodeString(tagNode))
	switch tagNode.pegRule {
	case ruleAuthor:
		p.setElement(Author, val)
	case ruleDataSource:
		p.setElement(DataSourceID, val)
	case ruleGenus:
		p.setElement(Genus, val)
	case ruleNameString:
		p.setElement(NameString, val)
	case ruleParentTaxon:
		p.setElement(ParentTaxon, val)
	case ruleSpecies:
		p.setElement(Species, val)
	case ruleSpeciesAny:
		p.setElement(SpeciesAny, val)
	case ruleSpeciesInfra:
		p.setElement(SpeciesInfra, val)
	case ruleUninomial:
		p.setElement(Uninomial, val)
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
