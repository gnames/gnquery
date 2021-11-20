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

func (p *parser) newElement(n *node32) {
	n = n.up
	tagNode := n.up
	valNode := n.next
	val := p.nodeString(valNode)

	fmt.Printf("TAG: %#v, %s\n", tagNode, p.nodeString(tagNode))
	switch tagNode.pegRule {
	case ruleAuthor:
		p.elements[Author] = val
	case ruleDataSource:
		p.elements[DataSource] = val
	case ruleGenus:
		p.elements[Genus] = val
	case ruleNameString:
		p.elements[NameString] = val
	case ruleParentTaxon:
		p.elements[ParentTaxon] = val
	case ruleSpecies:
		p.elements[Species] = val
	case ruleSpeciesAny:
		p.elements[SpeciesAny] = val
	case ruleSpeciesInfra:
		p.elements[SpeciesInfra] = val
	case ruleUninomial:
		p.elements[Uninomial] = val
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
