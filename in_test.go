package tagify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypeOf(t *testing.T) {
	assert.Equal(t, Unknown, ContentTypeOf("Unknown"))
	assert.Equal(t, HTML, ContentTypeOf("HTML"))
	assert.Equal(t, Text, ContentTypeOf("Text"))
}

func TestNewInFromString_ReadAllLines(t *testing.T) {
	in := NewInFromString("Test input reader of type text", Text)
	lns, err := in.ReadAllLines()
	assert.Nil(t, err)
	assert.Len(t, lns, 1)
}

func TestNewInFromString_ReadAllStrings(t *testing.T) {
	in := NewInFromString("Test input reader of type text", Text)
	strs, err := in.ReadAllStrings()
	assert.Nil(t, err)
	assert.Len(t, strs, 6)
}
