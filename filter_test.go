package tagify

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