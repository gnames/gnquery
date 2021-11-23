// package parser allows extract elements from a query string and create
// Input object required for GNames faceted search.
package parser

import "github.com/gnames/gnquery/ent/search"

// Parser contains methods for parsing a faceted search query for GNames.
type Parser interface {
	// ParseQuery takes a query string, parses it, registers warnings if
	// query is not well-formed, and packs received data into a search.Input
	// object.
	ParseQuery(q string) search.Input

	// Debug takes a query string, parses it, and prints aot an abstract
	// syntax tree of the parsed result. This method is used for developing
	// and debuging the parser behavior.
	Debug(string)
}
