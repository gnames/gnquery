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
	ruleTagAllResults
	ruleTagString
	ruleTagDS
	ruleTagYear
	ruleTagWord
	ruleAuthor
	ruleDataSource
	ruleAllResults
	ruleGenus
	ruleNameString
	ruleParentTaxon
	ruleSpecies
	ruleSpeciesAny
	ruleSpeciesInfra
	ruleYear
	ruleString
	ruleYearRange
	ruleYearNum
	ruleNumber
	ruleBool
	ruleDigits
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
	"TagAllResults",
	"TagString",
	"TagDS",
	"TagYear",
	"TagWord",
	"Author",
	"DataSource",
	"AllResults",
	"Genus",
	"NameString",
	"ParentTaxon",
	"Species",
	"SpeciesAny",
	"SpeciesInfra",
	"Year",
	"String",
	"YearRange",
	"YearNum",
	"Number",
	"Bool",
	"Digits",
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
	rules  [31]func() bool
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
		/* 3 Element <- <((TagAllResults Bool) / (TagString String) / (TagDS Number) / (TagYear (YearRange / YearNum)) / (TagWord Word))> */
		func() bool {
			position14, tokenIndex14 := position, tokenIndex
			{
				position15 := position
				{
					position16, tokenIndex16 := position, tokenIndex
					if !_rules[ruleTagAllResults]() {
						goto l17
					}
					if !_rules[ruleBool]() {
						goto l17
					}
					goto l16
				l17:
					position, tokenIndex = position16, tokenIndex16
					if !_rules[ruleTagString]() {
						goto l18
					}
					if !_rules[ruleString]() {
						goto l18
					}
					goto l16
				l18:
					position, tokenIndex = position16, tokenIndex16
					if !_rules[ruleTagDS]() {
						goto l19
					}
					if !_rules[ruleNumber]() {
						goto l19
					}
					goto l16
				l19:
					position, tokenIndex = position16, tokenIndex16
					if !_rules[ruleTagYear]() {
						goto l20
					}
					{
						position21, tokenIndex21 := position, tokenIndex
						if !_rules[ruleYearRange]() {
							goto l22
						}
						goto l21
					l22:
						position, tokenIndex = position21, tokenIndex21
						if !_rules[ruleYearNum]() {
							goto l20
						}
					}
				l21:
					goto l16
				l20:
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
		/* 4 TagAllResults <- <(AllResults ':')> */
		func() bool {
			position23, tokenIndex23 := position, tokenIndex
			{
				position24 := position
				if !_rules[ruleAllResults]() {
					goto l23
				}
				if buffer[position] != rune(':') {
					goto l23
				}
				position++
				add(ruleTagAllResults, position24)
			}
			return true
		l23:
			position, tokenIndex = position23, tokenIndex23
			return false
		},
		/* 5 TagString <- <(NameString ':')> */
		func() bool {
			position25, tokenIndex25 := position, tokenIndex
			{
				position26 := position
				if !_rules[ruleNameString]() {
					goto l25
				}
				if buffer[position] != rune(':') {
					goto l25
				}
				position++
				add(ruleTagString, position26)
			}
			return true
		l25:
			position, tokenIndex = position25, tokenIndex25
			return false
		},
		/* 6 TagDS <- <(DataSource ':')> */
		func() bool {
			position27, tokenIndex27 := position, tokenIndex
			{
				position28 := position
				if !_rules[ruleDataSource]() {
					goto l27
				}
				if buffer[position] != rune(':') {
					goto l27
				}
				position++
				add(ruleTagDS, position28)
			}
			return true
		l27:
			position, tokenIndex = position27, tokenIndex27
			return false
		},
		/* 7 TagYear <- <(Year ':')> */
		func() bool {
			position29, tokenIndex29 := position, tokenIndex
			{
				position30 := position
				if !_rules[ruleYear]() {
					goto l29
				}
				if buffer[position] != rune(':') {
					goto l29
				}
				position++
				add(ruleTagYear, position30)
			}
			return true
		l29:
			position, tokenIndex = position29, tokenIndex29
			return false
		},
		/* 8 TagWord <- <((Genus / ParentTaxon / SpeciesAny / Species / SpeciesInfra / Author) ':')> */
		func() bool {
			position31, tokenIndex31 := position, tokenIndex
			{
				position32 := position
				{
					position33, tokenIndex33 := position, tokenIndex
					if !_rules[ruleGenus]() {
						goto l34
					}
					goto l33
				l34:
					position, tokenIndex = position33, tokenIndex33
					if !_rules[ruleParentTaxon]() {
						goto l35
					}
					goto l33
				l35:
					position, tokenIndex = position33, tokenIndex33
					if !_rules[ruleSpeciesAny]() {
						goto l36
					}
					goto l33
				l36:
					position, tokenIndex = position33, tokenIndex33
					if !_rules[ruleSpecies]() {
						goto l37
					}
					goto l33
				l37:
					position, tokenIndex = position33, tokenIndex33
					if !_rules[ruleSpeciesInfra]() {
						goto l38
					}
					goto l33
				l38:
					position, tokenIndex = position33, tokenIndex33
					if !_rules[ruleAuthor]() {
						goto l31
					}
				}
			l33:
				if buffer[position] != rune(':') {
					goto l31
				}
				position++
				add(ruleTagWord, position32)
			}
			return true
		l31:
			position, tokenIndex = position31, tokenIndex31
			return false
		},
		/* 9 Author <- <(('a' 'u') / 'a')> */
		func() bool {
			position39, tokenIndex39 := position, tokenIndex
			{
				position40 := position
				{
					position41, tokenIndex41 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l42
					}
					position++
					if buffer[position] != rune('u') {
						goto l42
					}
					position++
					goto l41
				l42:
					position, tokenIndex = position41, tokenIndex41
					if buffer[position] != rune('a') {
						goto l39
					}
					position++
				}
			l41:
				add(ruleAuthor, position40)
			}
			return true
		l39:
			position, tokenIndex = position39, tokenIndex39
			return false
		},
		/* 10 DataSource <- <('d' 's')> */
		func() bool {
			position43, tokenIndex43 := position, tokenIndex
			{
				position44 := position
				if buffer[position] != rune('d') {
					goto l43
				}
				position++
				if buffer[position] != rune('s') {
					goto l43
				}
				position++
				add(ruleDataSource, position44)
			}
			return true
		l43:
			position, tokenIndex = position43, tokenIndex43
			return false
		},
		/* 11 AllResults <- <('a' 'l' 'l')> */
		func() bool {
			position45, tokenIndex45 := position, tokenIndex
			{
				position46 := position
				if buffer[position] != rune('a') {
					goto l45
				}
				position++
				if buffer[position] != rune('l') {
					goto l45
				}
				position++
				if buffer[position] != rune('l') {
					goto l45
				}
				position++
				add(ruleAllResults, position46)
			}
			return true
		l45:
			position, tokenIndex = position45, tokenIndex45
			return false
		},
		/* 12 Genus <- <(('g' 'e' 'n') / 'g')> */
		func() bool {
			position47, tokenIndex47 := position, tokenIndex
			{
				position48 := position
				{
					position49, tokenIndex49 := position, tokenIndex
					if buffer[position] != rune('g') {
						goto l50
					}
					position++
					if buffer[position] != rune('e') {
						goto l50
					}
					position++
					if buffer[position] != rune('n') {
						goto l50
					}
					position++
					goto l49
				l50:
					position, tokenIndex = position49, tokenIndex49
					if buffer[position] != rune('g') {
						goto l47
					}
					position++
				}
			l49:
				add(ruleGenus, position48)
			}
			return true
		l47:
			position, tokenIndex = position47, tokenIndex47
			return false
		},
		/* 13 NameString <- <'n'> */
		func() bool {
			position51, tokenIndex51 := position, tokenIndex
			{
				position52 := position
				if buffer[position] != rune('n') {
					goto l51
				}
				position++
				add(ruleNameString, position52)
			}
			return true
		l51:
			position, tokenIndex = position51, tokenIndex51
			return false
		},
		/* 14 ParentTaxon <- <('t' 'x')> */
		func() bool {
			position53, tokenIndex53 := position, tokenIndex
			{
				position54 := position
				if buffer[position] != rune('t') {
					goto l53
				}
				position++
				if buffer[position] != rune('x') {
					goto l53
				}
				position++
				add(ruleParentTaxon, position54)
			}
			return true
		l53:
			position, tokenIndex = position53, tokenIndex53
			return false
		},
		/* 15 Species <- <('s' 'p')> */
		func() bool {
			position55, tokenIndex55 := position, tokenIndex
			{
				position56 := position
				if buffer[position] != rune('s') {
					goto l55
				}
				position++
				if buffer[position] != rune('p') {
					goto l55
				}
				position++
				add(ruleSpecies, position56)
			}
			return true
		l55:
			position, tokenIndex = position55, tokenIndex55
			return false
		},
		/* 16 SpeciesAny <- <('s' 'p' '+')> */
		func() bool {
			position57, tokenIndex57 := position, tokenIndex
			{
				position58 := position
				if buffer[position] != rune('s') {
					goto l57
				}
				position++
				if buffer[position] != rune('p') {
					goto l57
				}
				position++
				if buffer[position] != rune('+') {
					goto l57
				}
				position++
				add(ruleSpeciesAny, position58)
			}
			return true
		l57:
			position, tokenIndex = position57, tokenIndex57
			return false
		},
		/* 17 SpeciesInfra <- <('i' 's' 'p')> */
		func() bool {
			position59, tokenIndex59 := position, tokenIndex
			{
				position60 := position
				if buffer[position] != rune('i') {
					goto l59
				}
				position++
				if buffer[position] != rune('s') {
					goto l59
				}
				position++
				if buffer[position] != rune('p') {
					goto l59
				}
				position++
				add(ruleSpeciesInfra, position60)
			}
			return true
		l59:
			position, tokenIndex = position59, tokenIndex59
			return false
		},
		/* 18 Year <- <(('y' 'r') / 'y')> */
		func() bool {
			position61, tokenIndex61 := position, tokenIndex
			{
				position62 := position
				{
					position63, tokenIndex63 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l64
					}
					position++
					if buffer[position] != rune('r') {
						goto l64
					}
					position++
					goto l63
				l64:
					position, tokenIndex = position63, tokenIndex63
					if buffer[position] != rune('y') {
						goto l61
					}
					position++
				}
			l63:
				add(ruleYear, position62)
			}
			return true
		l61:
			position, tokenIndex = position61, tokenIndex61
			return false
		},
		/* 19 String <- <(Word (_ Word)*)> */
		func() bool {
			position65, tokenIndex65 := position, tokenIndex
			{
				position66 := position
				if !_rules[ruleWord]() {
					goto l65
				}
			l67:
				{
					position68, tokenIndex68 := position, tokenIndex
					if !_rules[rule_]() {
						goto l68
					}
					if !_rules[ruleWord]() {
						goto l68
					}
					goto l67
				l68:
					position, tokenIndex = position68, tokenIndex68
				}
				add(ruleString, position66)
			}
			return true
		l65:
			position, tokenIndex = position65, tokenIndex65
			return false
		},
		/* 20 YearRange <- <(('-' YearNum) / (YearNum '-' YearNum?))> */
		func() bool {
			position69, tokenIndex69 := position, tokenIndex
			{
				position70 := position
				{
					position71, tokenIndex71 := position, tokenIndex
					if buffer[position] != rune('-') {
						goto l72
					}
					position++
					if !_rules[ruleYearNum]() {
						goto l72
					}
					goto l71
				l72:
					position, tokenIndex = position71, tokenIndex71
					if !_rules[ruleYearNum]() {
						goto l69
					}
					if buffer[position] != rune('-') {
						goto l69
					}
					position++
					{
						position73, tokenIndex73 := position, tokenIndex
						if !_rules[ruleYearNum]() {
							goto l73
						}
						goto l74
					l73:
						position, tokenIndex = position73, tokenIndex73
					}
				l74:
				}
			l71:
				add(ruleYearRange, position70)
			}
			return true
		l69:
			position, tokenIndex = position69, tokenIndex69
			return false
		},
		/* 21 YearNum <- <(('1' / '2') ('0' / '7' / '8' / '9') Digits Digits)> */
		func() bool {
			position75, tokenIndex75 := position, tokenIndex
			{
				position76 := position
				{
					position77, tokenIndex77 := position, tokenIndex
					if buffer[position] != rune('1') {
						goto l78
					}
					position++
					goto l77
				l78:
					position, tokenIndex = position77, tokenIndex77
					if buffer[position] != rune('2') {
						goto l75
					}
					position++
				}
			l77:
				{
					position79, tokenIndex79 := position, tokenIndex
					if buffer[position] != rune('0') {
						goto l80
					}
					position++
					goto l79
				l80:
					position, tokenIndex = position79, tokenIndex79
					if buffer[position] != rune('7') {
						goto l81
					}
					position++
					goto l79
				l81:
					position, tokenIndex = position79, tokenIndex79
					if buffer[position] != rune('8') {
						goto l82
					}
					position++
					goto l79
				l82:
					position, tokenIndex = position79, tokenIndex79
					if buffer[position] != rune('9') {
						goto l75
					}
					position++
				}
			l79:
				if !_rules[ruleDigits]() {
					goto l75
				}
				if !_rules[ruleDigits]() {
					goto l75
				}
				add(ruleYearNum, position76)
			}
			return true
		l75:
			position, tokenIndex = position75, tokenIndex75
			return false
		},
		/* 22 Number <- <[0-9]+> */
		func() bool {
			position83, tokenIndex83 := position, tokenIndex
			{
				position84 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l83
				}
				position++
			l85:
				{
					position86, tokenIndex86 := position, tokenIndex
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l86
					}
					position++
					goto l85
				l86:
					position, tokenIndex = position86, tokenIndex86
				}
				add(ruleNumber, position84)
			}
			return true
		l83:
			position, tokenIndex = position83, tokenIndex83
			return false
		},
		/* 23 Bool <- <(('t' 'r' 'u' 'e') / 't' / ('f' 'a' 'l' 's' 'e') / 'f')> */
		func() bool {
			position87, tokenIndex87 := position, tokenIndex
			{
				position88 := position
				{
					position89, tokenIndex89 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l90
					}
					position++
					if buffer[position] != rune('r') {
						goto l90
					}
					position++
					if buffer[position] != rune('u') {
						goto l90
					}
					position++
					if buffer[position] != rune('e') {
						goto l90
					}
					position++
					goto l89
				l90:
					position, tokenIndex = position89, tokenIndex89
					if buffer[position] != rune('t') {
						goto l91
					}
					position++
					goto l89
				l91:
					position, tokenIndex = position89, tokenIndex89
					if buffer[position] != rune('f') {
						goto l92
					}
					position++
					if buffer[position] != rune('a') {
						goto l92
					}
					position++
					if buffer[position] != rune('l') {
						goto l92
					}
					position++
					if buffer[position] != rune('s') {
						goto l92
					}
					position++
					if buffer[position] != rune('e') {
						goto l92
					}
					position++
					goto l89
				l92:
					position, tokenIndex = position89, tokenIndex89
					if buffer[position] != rune('f') {
						goto l87
					}
					position++
				}
			l89:
				add(ruleBool, position88)
			}
			return true
		l87:
			position, tokenIndex = position87, tokenIndex87
			return false
		},
		/* 24 Digits <- <[0-9]> */
		func() bool {
			position93, tokenIndex93 := position, tokenIndex
			{
				position94 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l93
				}
				position++
				add(ruleDigits, position94)
			}
			return true
		l93:
			position, tokenIndex = position93, tokenIndex93
			return false
		},
		/* 25 Word <- <((!(':' / ' ') .)+ &(_ / END))> */
		func() bool {
			position95, tokenIndex95 := position, tokenIndex
			{
				position96 := position
				{
					position99, tokenIndex99 := position, tokenIndex
					{
						position100, tokenIndex100 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l101
						}
						position++
						goto l100
					l101:
						position, tokenIndex = position100, tokenIndex100
						if buffer[position] != rune(' ') {
							goto l99
						}
						position++
					}
				l100:
					goto l95
				l99:
					position, tokenIndex = position99, tokenIndex99
				}
				if !matchDot() {
					goto l95
				}
			l97:
				{
					position98, tokenIndex98 := position, tokenIndex
					{
						position102, tokenIndex102 := position, tokenIndex
						{
							position103, tokenIndex103 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l104
							}
							position++
							goto l103
						l104:
							position, tokenIndex = position103, tokenIndex103
							if buffer[position] != rune(' ') {
								goto l102
							}
							position++
						}
					l103:
						goto l98
					l102:
						position, tokenIndex = position102, tokenIndex102
					}
					if !matchDot() {
						goto l98
					}
					goto l97
				l98:
					position, tokenIndex = position98, tokenIndex98
				}
				{
					position105, tokenIndex105 := position, tokenIndex
					{
						position106, tokenIndex106 := position, tokenIndex
						if !_rules[rule_]() {
							goto l107
						}
						goto l106
					l107:
						position, tokenIndex = position106, tokenIndex106
						if !_rules[ruleEND]() {
							goto l95
						}
					}
				l106:
					position, tokenIndex = position105, tokenIndex105
				}
				add(ruleWord, position96)
			}
			return true
		l95:
			position, tokenIndex = position95, tokenIndex95
			return false
		},
		/* 26 _ <- <(MultipleSpace / SingleSpace)> */
		func() bool {
			position108, tokenIndex108 := position, tokenIndex
			{
				position109 := position
				{
					position110, tokenIndex110 := position, tokenIndex
					if !_rules[ruleMultipleSpace]() {
						goto l111
					}
					goto l110
				l111:
					position, tokenIndex = position110, tokenIndex110
					if !_rules[ruleSingleSpace]() {
						goto l108
					}
				}
			l110:
				add(rule_, position109)
			}
			return true
		l108:
			position, tokenIndex = position108, tokenIndex108
			return false
		},
		/* 27 MultipleSpace <- <(SingleSpace SingleSpace+)> */
		func() bool {
			position112, tokenIndex112 := position, tokenIndex
			{
				position113 := position
				if !_rules[ruleSingleSpace]() {
					goto l112
				}
				if !_rules[ruleSingleSpace]() {
					goto l112
				}
			l114:
				{
					position115, tokenIndex115 := position, tokenIndex
					if !_rules[ruleSingleSpace]() {
						goto l115
					}
					goto l114
				l115:
					position, tokenIndex = position115, tokenIndex115
				}
				add(ruleMultipleSpace, position113)
			}
			return true
		l112:
			position, tokenIndex = position112, tokenIndex112
			return false
		},
		/* 28 SingleSpace <- <' '> */
		func() bool {
			position116, tokenIndex116 := position, tokenIndex
			{
				position117 := position
				if buffer[position] != rune(' ') {
					goto l116
				}
				position++
				add(ruleSingleSpace, position117)
			}
			return true
		l116:
			position, tokenIndex = position116, tokenIndex116
			return false
		},
		/* 29 END <- <!.> */
		func() bool {
			position118, tokenIndex118 := position, tokenIndex
			{
				position119 := position
				{
					position120, tokenIndex120 := position, tokenIndex
					if !matchDot() {
						goto l120
					}
					goto l118
				l120:
					position, tokenIndex = position120, tokenIndex120
				}
				add(ruleEND, position119)
			}
			return true
		l118:
			position, tokenIndex = position118, tokenIndex118
			return false
		},
	}
	p.rules = _rules
	return nil
}
