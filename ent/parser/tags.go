package parser

type Tag int

const (
	Author Tag = iota
	DataSourceID
	ParentTaxon
	Year
	NameString
	Uninomial
	Genus
	Species
	SpeciesAny
	SpeciesInfra
)

var tagString = map[Tag]string{
	Author:       "au",
	DataSourceID: "ds",
	ParentTaxon:  "tx",
	Year:         "y",
	NameString:   "n",
	Uninomial:    "u",
	Genus:        "g",
	Species:      "sp",
	SpeciesAny:   "sp+",
	SpeciesInfra: "isp",
}

func (t Tag) String() string {
	return tagString[t]
}
