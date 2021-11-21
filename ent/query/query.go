package query

type Query struct {
	Input        string   `json:"input"`
	Warnings     []string `json:"warnings,omitempty"`
	DataSourceID int      `json:"dataSourceID,omitempty"`
	ParentTaxon  string   `json:"parentTaxon,omitempty"`
	NameString   string   `json:"nameString,omitempty"`
	Uninomial    string   `json:"uninomial,omitempty"`
	Genus        string   `json:"genus,omitempty"`
	SpeciesAny   string   `json:"speciesAny,omitempty"`
	Species      string   `json:"species,omitempty"`
	SpeciesInfra string   `json:"speciesInfra,omitempty"`
	Author       string   `json:"author,omitempty"`
	Year         int      `json:"year,omitempty"`
	Tail         string   `json:"tail,omitempty"`
}
