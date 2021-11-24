package parser

// Code generated by peg query.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleQuery
	ruleTail
	ruleComponents
	ruleElement
	ruleTagString
	ruleTagNum
	ruleTagWord
	ruleAuthor
	ruleDataSource
	ruleGenus
	ruleNameString
	ruleParentTaxon
	ruleSpecies
	ruleSpeciesAny
	ruleSpeciesInfra
	ruleUninomial
	ruleYear
	ruleString
	ruleNumber
	ruleWord
	rule_
	ruleMultipleSpace
	ruleSingleSpace
	ruleEND
)

var rul3s = [...]string{
	"Unknown",
	"Query",
	"Tail",
	"Components",
	"Element",
	"TagString",
	"TagNum",
	"TagWord",
	"Author",
	"DataSource",
	"Genus",
	"NameString",
	"ParentTaxon",
	"Species",
	"SpeciesAny",
	"SpeciesInfra",
	"Uninomial",
	"Year",
	"String",
	"Number",
	"Word",
	"_",
	"MultipleSpace",
	"SingleSpace",
	"END",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type engine struct {
	Buffer string
	buffer []rune
	rules  [25]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *engine) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *engine) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *engine
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *engine) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *engine) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *engine) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func Pretty(pretty bool) func(*engine) error {
	return func(p *engine) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*engine) error {
	return func(p *engine) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *engine) Init(options ...func(*engine) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Query <- <(_? Components Tail END)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				{
					position2, tokenIndex2 := position, tokenIndex
					if !_rules[rule_]() {
						goto l2
					}
					goto l3
				l2:
					position, tokenIndex = position2, tokenIndex2
				}
			l3:
				if !_rules[ruleComponents]() {
					goto l0
				}
				if !_rules[ruleTail]() {
					goto l0
				}
				if !_rules[ruleEND]() {
					goto l0
				}
				add(ruleQuery, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Tail <- <(_ .*)?> */
		func() bool {
			{
				position5 := position
				{
					position6, tokenIndex6 := position, tokenIndex
					if !_rules[rule_]() {
						goto l6
					}
				l8:
					{
						position9, tokenIndex9 := position, tokenIndex
						if !matchDot() {
							goto l9
						}
						goto l8
					l9:
						position, tokenIndex = position9, tokenIndex9
					}
					goto l7
				l6:
					position, tokenIndex = position6, tokenIndex6
				}
			l7:
				add(ruleTail, position5)
			}
			return true
		},
		/* 2 Components <- <(Element (_ Element)*)> */
		func() bool {
			position10, tokenIndex10 := position, tokenIndex
			{
				position11 := position
				if !_rules[ruleElement]() {
					goto l10
				}
			l12:
				{
					position13, tokenIndex13 := position, tokenIndex
					if !_rules[rule_]() {
						goto l13
					}
					if !_rules[ruleElement]() {
						goto l13
					}
					goto l12
				l13:
					position, tokenIndex = position13, tokenIndex13
				}
				add(ruleComponents, position11)
			}
			return true
		l10:
			position, tokenIndex = position10, tokenIndex10
			return false
		},
		/* 3 Element <- <((TagString String) / (TagNum Number) / (TagWord Word))> */
		func() bool {
			position14, tokenIndex14 := position, tokenIndex
			{
				position15 := position
				{
					position16, tokenIndex16 := position, tokenIndex
					if !_rules[ruleTagString]() {
						goto l17
					}
					if !_rules[ruleString]() {
						goto l17
					}
					goto l16
				l17:
					position, tokenIndex = position16, tokenIndex16
					if !_rules[ruleTagNum]() {
						goto l18
					}
					if !_rules[ruleNumber]() {
						goto l18
					}
					goto l16
				l18:
					position, tokenIndex = position16, tokenIndex16
					if !_rules[ruleTagWord]() {
						goto l14
					}
					if !_rules[ruleWord]() {
						goto l14
					}
				}
			l16:
				add(ruleElement, position15)
			}
			return true
		l14:
			position, tokenIndex = position14, tokenIndex14
			return false
		},
		/* 4 TagString <- <(NameString ':')> */
		func() bool {
			position19, tokenIndex19 := position, tokenIndex
			{
				position20 := position
				if !_rules[ruleNameString]() {
					goto l19
				}
				if buffer[position] != rune(':') {
					goto l19
				}
				position++
				add(ruleTagString, position20)
			}
			return true
		l19:
			position, tokenIndex = position19, tokenIndex19
			return false
		},
		/* 5 TagNum <- <((DataSource / Year) ':')> */
		func() bool {
			position21, tokenIndex21 := position, tokenIndex
			{
				position22 := position
				{
					position23, tokenIndex23 := position, tokenIndex
					if !_rules[ruleDataSource]() {
						goto l24
					}
					goto l23
				l24:
					position, tokenIndex = position23, tokenIndex23
					if !_rules[ruleYear]() {
						goto l21
					}
				}
			l23:
				if buffer[position] != rune(':') {
					goto l21
				}
				position++
				add(ruleTagNum, position22)
			}
			return true
		l21:
			position, tokenIndex = position21, tokenIndex21
			return false
		},
		/* 6 TagWord <- <((Genus / ParentTaxon / SpeciesAny / Species / SpeciesInfra / Uninomial / Author) ':')> */
		func() bool {
			position25, tokenIndex25 := position, tokenIndex
			{
				position26 := position
				{
					position27, tokenIndex27 := position, tokenIndex
					if !_rules[ruleGenus]() {
						goto l28
					}
					goto l27
				l28:
					position, tokenIndex = position27, tokenIndex27
					if !_rules[ruleParentTaxon]() {
						goto l29
					}
					goto l27
				l29:
					position, tokenIndex = position27, tokenIndex27
					if !_rules[ruleSpeciesAny]() {
						goto l30
					}
					goto l27
				l30:
					position, tokenIndex = position27, tokenIndex27
					if !_rules[ruleSpecies]() {
						goto l31
					}
					goto l27
				l31:
					position, tokenIndex = position27, tokenIndex27
					if !_rules[ruleSpeciesInfra]() {
						goto l32
					}
					goto l27
				l32:
					position, tokenIndex = position27, tokenIndex27
					if !_rules[ruleUninomial]() {
						goto l33
					}
					goto l27
				l33:
					position, tokenIndex = position27, tokenIndex27
					if !_rules[ruleAuthor]() {
						goto l25
					}
				}
			l27:
				if buffer[position] != rune(':') {
					goto l25
				}
				position++
				add(ruleTagWord, position26)
			}
			return true
		l25:
			position, tokenIndex = position25, tokenIndex25
			return false
		},
		/* 7 Author <- <(('a' 'u') / 'a')> */
		func() bool {
			position34, tokenIndex34 := position, tokenIndex
			{
				position35 := position
				{
					position36, tokenIndex36 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l37
					}
					position++
					if buffer[position] != rune('u') {
						goto l37
					}
					position++
					goto l36
				l37:
					position, tokenIndex = position36, tokenIndex36
					if buffer[position] != rune('a') {
						goto l34
					}
					position++
				}
			l36:
				add(ruleAuthor, position35)
			}
			return true
		l34:
			position, tokenIndex = position34, tokenIndex34
			return false
		},
		/* 8 DataSource <- <('d' 's')> */
		func() bool {
			position38, tokenIndex38 := position, tokenIndex
			{
				position39 := position
				if buffer[position] != rune('d') {
					goto l38
				}
				position++
				if buffer[position] != rune('s') {
					goto l38
				}
				position++
				add(ruleDataSource, position39)
			}
			return true
		l38:
			position, tokenIndex = position38, tokenIndex38
			return false
		},
		/* 9 Genus <- <(('g' 'e' 'n') / 'g')> */
		func() bool {
			position40, tokenIndex40 := position, tokenIndex
			{
				position41 := position
				{
					position42, tokenIndex42 := position, tokenIndex
					if buffer[position] != rune('g') {
						goto l43
					}
					position++
					if buffer[position] != rune('e') {
						goto l43
					}
					position++
					if buffer[position] != rune('n') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex = position42, tokenIndex42
					if buffer[position] != rune('g') {
						goto l40
					}
					position++
				}
			l42:
				add(ruleGenus, position41)
			}
			return true
		l40:
			position, tokenIndex = position40, tokenIndex40
			return false
		},
		/* 10 NameString <- <'n'> */
		func() bool {
			position44, tokenIndex44 := position, tokenIndex
			{
				position45 := position
				if buffer[position] != rune('n') {
					goto l44
				}
				position++
				add(ruleNameString, position45)
			}
			return true
		l44:
			position, tokenIndex = position44, tokenIndex44
			return false
		},
		/* 11 ParentTaxon <- <('t' 'x')> */
		func() bool {
			position46, tokenIndex46 := position, tokenIndex
			{
				position47 := position
				if buffer[position] != rune('t') {
					goto l46
				}
				position++
				if buffer[position] != rune('x') {
					goto l46
				}
				position++
				add(ruleParentTaxon, position47)
			}
			return true
		l46:
			position, tokenIndex = position46, tokenIndex46
			return false
		},
		/* 12 Species <- <('s' 'p')> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				if buffer[position] != rune('s') {
					goto l48
				}
				position++
				if buffer[position] != rune('p') {
					goto l48
				}
				position++
				add(ruleSpecies, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 13 SpeciesAny <- <('s' 'p' '+')> */
		func() bool {
			position50, tokenIndex50 := position, tokenIndex
			{
				position51 := position
				if buffer[position] != rune('s') {
					goto l50
				}
				position++
				if buffer[position] != rune('p') {
					goto l50
				}
				position++
				if buffer[position] != rune('+') {
					goto l50
				}
				position++
				add(ruleSpeciesAny, position51)
			}
			return true
		l50:
			position, tokenIndex = position50, tokenIndex50
			return false
		},
		/* 14 SpeciesInfra <- <('i' 's' 'p')> */
		func() bool {
			position52, tokenIndex52 := position, tokenIndex
			{
				position53 := position
				if buffer[position] != rune('i') {
					goto l52
				}
				position++
				if buffer[position] != rune('s') {
					goto l52
				}
				position++
				if buffer[position] != rune('p') {
					goto l52
				}
				position++
				add(ruleSpeciesInfra, position53)
			}
			return true
		l52:
			position, tokenIndex = position52, tokenIndex52
			return false
		},
		/* 15 Uninomial <- <'u'> */
		func() bool {
			position54, tokenIndex54 := position, tokenIndex
			{
				position55 := position
				if buffer[position] != rune('u') {
					goto l54
				}
				position++
				add(ruleUninomial, position55)
			}
			return true
		l54:
			position, tokenIndex = position54, tokenIndex54
			return false
		},
		/* 16 Year <- <(('y' 'r') / 'y')> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				{
					position58, tokenIndex58 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l59
					}
					position++
					if buffer[position] != rune('r') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex = position58, tokenIndex58
					if buffer[position] != rune('y') {
						goto l56
					}
					position++
				}
			l58:
				add(ruleYear, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 17 String <- <(Word (_ Word)*)> */
		func() bool {
			position60, tokenIndex60 := position, tokenIndex
			{
				position61 := position
				if !_rules[ruleWord]() {
					goto l60
				}
			l62:
				{
					position63, tokenIndex63 := position, tokenIndex
					if !_rules[rule_]() {
						goto l63
					}
					if !_rules[ruleWord]() {
						goto l63
					}
					goto l62
				l63:
					position, tokenIndex = position63, tokenIndex63
				}
				add(ruleString, position61)
			}
			return true
		l60:
			position, tokenIndex = position60, tokenIndex60
			return false
		},
		/* 18 Number <- <[0-9]+> */
		func() bool {
			position64, tokenIndex64 := position, tokenIndex
			{
				position65 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l64
				}
				position++
			l66:
				{
					position67, tokenIndex67 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l67
					}
					position++
					goto l66
				l67:
					position, tokenIndex = position67, tokenIndex67
				}
				add(ruleNumber, position65)
			}
			return true
		l64:
			position, tokenIndex = position64, tokenIndex64
			return false
		},
		/* 19 Word <- <((!(':' / ' ') .)+ &(_ / END))> */
		func() bool {
			position68, tokenIndex68 := position, tokenIndex
			{
				position69 := position
				{
					position72, tokenIndex72 := position, tokenIndex
					{
						position73, tokenIndex73 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l74
						}
						position++
						goto l73
					l74:
						position, tokenIndex = position73, tokenIndex73
						if buffer[position] != rune(' ') {
							goto l72
						}
						position++
					}
				l73:
					goto l68
				l72:
					position, tokenIndex = position72, tokenIndex72
				}
				if !matchDot() {
					goto l68
				}
			l70:
				{
					position71, tokenIndex71 := position, tokenIndex
					{
						position75, tokenIndex75 := position, tokenIndex
						{
							position76, tokenIndex76 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l77
							}
							position++
							goto l76
						l77:
							position, tokenIndex = position76, tokenIndex76
							if buffer[position] != rune(' ') {
								goto l75
							}
							position++
						}
					l76:
						goto l71
					l75:
						position, tokenIndex = position75, tokenIndex75
					}
					if !matchDot() {
						goto l71
					}
					goto l70
				l71:
					position, tokenIndex = position71, tokenIndex71
				}
				{
					position78, tokenIndex78 := position, tokenIndex
					{
						position79, tokenIndex79 := position, tokenIndex
						if !_rules[rule_]() {
							goto l80
						}
						goto l79
					l80:
						position, tokenIndex = position79, tokenIndex79
						if !_rules[ruleEND]() {
							goto l68
						}
					}
				l79:
					position, tokenIndex = position78, tokenIndex78
				}
				add(ruleWord, position69)
			}
			return true
		l68:
			position, tokenIndex = position68, tokenIndex68
			return false
		},
		/* 20 _ <- <(MultipleSpace / SingleSpace)> */
		func() bool {
			position81, tokenIndex81 := position, tokenIndex
			{
				position82 := position
				{
					position83, tokenIndex83 := position, tokenIndex
					if !_rules[ruleMultipleSpace]() {
						goto l84
					}
					goto l83
				l84:
					position, tokenIndex = position83, tokenIndex83
					if !_rules[ruleSingleSpace]() {
						goto l81
					}
				}
			l83:
				add(rule_, position82)
			}
			return true
		l81:
			position, tokenIndex = position81, tokenIndex81
			return false
		},
		/* 21 MultipleSpace <- <(SingleSpace SingleSpace+)> */
		func() bool {
			position85, tokenIndex85 := position, tokenIndex
			{
				position86 := position
				if !_rules[ruleSingleSpace]() {
					goto l85
				}
				if !_rules[ruleSingleSpace]() {
					goto l85
				}
			l87:
				{
					position88, tokenIndex88 := position, tokenIndex
					if !_rules[ruleSingleSpace]() {
						goto l88
					}
					goto l87
				l88:
					position, tokenIndex = position88, tokenIndex88
				}
				add(ruleMultipleSpace, position86)
			}
			return true
		l85:
			position, tokenIndex = position85, tokenIndex85
			return false
		},
		/* 22 SingleSpace <- <' '> */
		func() bool {
			position89, tokenIndex89 := position, tokenIndex
			{
				position90 := position
				if buffer[position] != rune(' ') {
					goto l89
				}
				position++
				add(ruleSingleSpace, position90)
			}
			return true
		l89:
			position, tokenIndex = position89, tokenIndex89
			return false
		},
		/* 23 END <- <!.> */
		func() bool {
			position91, tokenIndex91 := position, tokenIndex
			{
				position92 := position
				{
					position93, tokenIndex93 := position, tokenIndex
					if !matchDot() {
						goto l93
					}
					goto l91
				l93:
					position, tokenIndex = position93, tokenIndex93
				}
				add(ruleEND, position92)
			}
			return true
		l91:
			position, tokenIndex = position91, tokenIndex91
			return false
		},
	}
	p.rules = _rules
	return nil
}
