package parser_test

import (
	"testing"

	"github.com/gnames/gnquery/ent/parser"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	q := "n:Bubo bubo au:Linn. sp:caboom all:t ds:1,2,3"
	parser.New().Debug(q)
}

func TestParseQuery(t *testing.T) {
	s := "n:Bubo bubo au:Linn. sp:caboom all:t abr:t ds:1,2"
	p := parser.New()
	q := p.ParseQuery(s)
	assert.True(t, len(q.Warnings) > 0, "warn")
	assert.Equal(t, "n:Bubo bubo au:Linn. sp:caboom all:t abr:t ds:1,2", q.Query, "query")
	assert.Equal(t, "Bubo bubo", q.NameString, "n")
	assert.Equal(t, "bubo", q.Species, "sp")
	assert.Equal(t, "", q.ParentTaxon, "tx")
	assert.Equal(t, []int{1, 2}, q.DataSources, "ds")
	assert.Equal(t, true, q.WithAllMatches, "all")
	assert.Equal(t, "Linn.", q.Author, "au")
	assert.Equal(t, true, q.WithAllBestResults, "abr")
}

func TestNameString(t *testing.T) {
	tests := []struct {
		msg, q   string
		warns    int
		gen      string
		sp       string
		isp      string
		au       string
		yr       int
		hasRange bool
		yst      int
		yend     int
		tail     string
	}{
		{"normal", "n:Bubo bubo Linn. 1756", 0, "Bubo", "bubo", "", "Linn.", 1756, false, 0, 0, ""},
		{
			"dupl sp",
			"n:Bubo bubo Linn. 1756 sp:alba",
			1,
			"Bubo",
			"bubo",
			"",
			"Linn.",
			1756,
			false,
			0,
			0,
			"",
		},
		{
			"dupl gen sp",
			"n:Bubo bubo Linn. 1756 g:Betula sp:alba",
			2,
			"Bubo",
			"bubo",
			"",
			"Linn.",
			1756,
			false,
			0,
			0,
			"",
		},
		{
			"au dupl",
			"n:Bubo bubo Linn. 1756 au:Kenth.",
			1,
			"Bubo",
			"bubo",
			"",
			"Linn.",
			1756,
			false,
			0,
			0,
			"",
		},
		{"au 1char", "n:Bubo bubo L 1756", 0, "Bubo", "bubo", "", "L", 1756, false, 0, 0, ""},
		{
			"no au",
			"n:Bubo bubo 1756 au:Kenth.",
			0,
			"Bubo",
			"bubo",
			"",
			"Kenth.",
			1756,
			false,
			0,
			0,
			"",
		},
		{
			"bad ord",
			"n:Bubo bubo 1756 Linn. au:Kenth.",
			1,
			"Bubo",
			"bubo",
			"",
			"",
			1756,
			false,
			0,
			0,
			"",
		},
		{"isp", "n:Bubo alba bubo Linn.", 0, "Bubo", "alba", "bubo", "Linn.", 0, false, 0, 0, ""},
		{
			"range",
			"n:Bubo bubo Linn. 1756-1777",
			0,
			"Bubo",
			"bubo",
			"",
			"Linn.",
			0,
			true,
			1756,
			1777,
			"",
		},
		{"range2", "n:Bubo bubo Linn. -1777", 0, "Bubo", "bubo", "", "Linn.", 0, true, 0, 1777, ""},
		{"range3", "n:Bubo bubo Linn. 1888-", 0, "Bubo", "bubo", "", "Linn.", 0, true, 1888, 0, ""},
	}

	p := parser.New()
	for _, v := range tests {
		res := p.ParseQuery(v.q)
		assert.Equal(t, v.q, res.Query, v.msg)
		assert.True(t, len(res.Warnings) == v.warns, v.msg)
		assert.Equal(t, v.gen, res.Genus, v.msg)
		assert.Equal(t, v.sp, res.Species, v.msg)
		assert.Equal(t, v.isp, res.SpeciesInfra, v.msg)
		assert.Equal(t, v.au, res.Author, v.msg)
		assert.Equal(t, v.yr, res.Year, v.msg)
		if v.hasRange {
			assert.Equal(t, v.yst, res.YearStart)
			assert.Equal(t, v.yend, res.YearEnd)
		}
	}
}

func TestQueries(t *testing.T) {
	tests := []struct {
		msg, q string
		warns  int
		yst    int
		yend   int
		ds     []int
		all    bool
		abr    bool
	}{
		{"full range", "g:B. sp:b. y:1888-2000", 0, 1888, 2000, nil, false, false},
		{"range", "g:B. sp:b. y:1888-2000", 0, 1888, 2000, nil, false, false},
		{"greater", "g:B. sp:b. y:1888-", 0, 1888, 0, nil, false, false},
		{"less", "g:B. sp:b. y:-2000", 0, 0, 2000, nil, false, false},
		{"all", "g:B. sp:b. y:-2000 all:t", 0, 0, 2000, nil, true, false},
		{"abr", "g:B. sp:b. y:-2000 all:t abr:t", 0, 0, 2000, nil, true, true},
		{"all2", "g:B. sp:b. y:-2000 all:true", 0, 0, 2000, nil, true, false},
		{"best only", "g:B. sp:b. y:-2000 all:f", 0, 0, 2000, nil, false, false},
		{"best only", "g:B. sp:b. y:-2000 all:false", 0, 0, 2000, nil, false, false},
		{
			"mult ds",
			"g:B. sp:b. ds:1,2,3 y:-2000 all:false",
			0,
			0,
			2000,
			[]int{1, 2, 3},
			false,
			false,
		},
		{"single ds", "g:B. sp:b. ds:1 y:-2000 all:false", 0, 0, 2000, []int{1}, false, false},
		{
			"single ds",
			"g:B. sp:b. ds:0 y:-2000 all:false abr:true",
			0,
			0,
			2000,
			[]int{0},
			false,
			true,
		},
		{"single ds", "g:B. sp:b. ds:0 y:-2000 all:t", 0, 0, 2000, []int{0}, true, false},
	}

	p := parser.New()

	for _, v := range tests {
		res := p.ParseQuery(v.q)
		assert.True(t, len(res.Warnings) == v.warns)
		assert.Equal(t, v.all, res.WithAllMatches)
		assert.Equal(t, v.abr, res.WithAllBestResults)
		assert.Equal(t, 0, res.Year)
		assert.NotNil(t, res.YearRange)
		assert.Equal(t, v.yst, res.YearStart)
		assert.Equal(t, v.yend, res.YearEnd)
		assert.Equal(t, v.ds, res.DataSources)
	}

	q := "g:B. sp:b. y:-"
	res := p.ParseQuery(q)
	assert.True(t, len(res.Warnings) > 0)
	assert.Nil(t, res.YearRange)
}
