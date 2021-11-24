package tag_test

import (
	"testing"

	"github.com/gnames/gnquery/ent/tag"
	"github.com/stretchr/testify/assert"
)

func TestTagString(t *testing.T) {
	res := tag.SpeciesAny.String()
	assert.Equal(t, res, "sp+")
}
