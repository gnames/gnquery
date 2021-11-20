package parser

type Tag int

const (
	Author Tag = iota
	DataSource
	Genus
	NameString
	ParentTaxon
	Species
	SpeciesAny
	SpeciesInfra
	Year
	Uninomial
)
