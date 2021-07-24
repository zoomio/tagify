package tagify

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zoomio/tagify/config"
)

func TestContentTypeOf(t *testing.T) {
	assert.Equal(t, config.Unknown, config.ContentTypeOf("Unknown"))
	assert.Equal(t, config.HTML, config.ContentTypeOf("HTML"))
	assert.Equal(t, config.Text, config.ContentTypeOf("Text"))
}

func TestNewInFromString_ReadLines(t *testing.T) {
	in := newInFromString("Test input reader of type text", config.Text)
	lns, err := in.ReadLines()
	assert.Nil(t, err)
	assert.Len(t, lns, 1)
}
