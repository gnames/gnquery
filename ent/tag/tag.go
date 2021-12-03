package tag

// Tag represents allowed tags in a search query. The tags in the query
// should be followed by a colon character. Only NameString tag is allowed
// to have more than one word.
//
// For example `g:Hom. sp:sapiens`, `n:Homo sapiens`
type Tag int

const (
	AllResults Tag = iota
	Author
	DataSourceIDs
	ParentTaxon
	Year
	NameString
	Genus
	Species
	SpeciesAny
	SpeciesInfra
)

var tagString = map[Tag]string{
	AllResults:    "all",
	Author:        "au",
	DataSourceIDs: "ds",
	ParentTaxon:   "tx",
	Year:          "y",
	NameString:    "n",
	Genus:         "g",
	Species:       "sp",
	SpeciesAny:    "asp",
	SpeciesInfra:  "isp",
}

// String returns a string representation of a tag.
func (t Tag) String() string {
	return tagString[t]
}
