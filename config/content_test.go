package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypeOf(t *testing.T) {
	assert.Equal(t, Unknown, ContentTypeOf("Unknown"))
	assert.Equal(t, HTML, ContentTypeOf("HTML"))
	assert.Equal(t, Text, ContentTypeOf("Text"))
}
