package rank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	assert.Equal(t, "part", sanitize("(part)"))
	assert.Equal(t, "part", sanitize("(part1)"))
	assert.Equal(t, "part", sanitize("part1"))
	assert.Equal(t, "part", sanitize("part}1"))
	assert.Equal(t, "part", sanitize("{part}"))
	assert.Equal(t, "d'arko", sanitize("{d'arko}"))
	assert.Equal(t, "", sanitize("1)"))
}

func TestToSingular(t *testing.T) {
	assert.Equal(t, "part", toSingular("parts"))
	assert.Equal(t, "algorithm", toSingular("algorithms"))
	assert.Equal(t, "year", toSingular("years"))
	assert.Equal(t, "cat", toSingular("cats"))
	assert.Equal(t, "person", toSingular("people"))
}

func TestDeDupe(t *testing.T) {
	
}