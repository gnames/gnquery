package parser

// Debug takes a string, parses it, and prints its AST.
func (p *Engine) Debug(s string) {
	p.PrettyPrintSyntaxTree(s)
}
