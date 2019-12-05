package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoomio/inout"
)

var text = `
	There was a boy
	Who's name was Jim
`

func Test_ParseText_Empty(t *testing.T) {
	tags := ParseText(inout.NewFromString(""), false, false)
	assert.Len(t, tags, 0)
}

func Test_ParseText_WithStopWords(t *testing.T) {
	tags := ParseText(inout.NewFromString(text), false, false)
	assert.Len(t, tags, 7)
	assert.Subset(t, ToStrings(tags), []string{"there", "was", "a", "boy", "who's", "name", "jim"})
}

func Test_ParseText_NoStopWords(t *testing.T) {
	tags := ParseText(inout.NewFromString(text), false, true)
	assert.Len(t, tags, 2)
	assert.Subset(t, ToStrings(tags), []string{"boy", "jim"})
}
