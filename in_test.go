package tagify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInFromString_ReadLines(t *testing.T) {
	in := newInFromString("Test input reader of type text", Text)
	lns, err := in.ReadLines()
	assert.Nil(t, err)
	assert.Len(t, lns, 1)
}
