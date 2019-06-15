package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var html = []string{
	"<html>",
	"<body>",
	"<p>There was a boy</p>",
	"<p>Who's name was Jim.</p>",
	"</body>",
	"</html>",
}

func Test_ParseHTML_Empty(t *testing.T) {
	tags := ParseHTML([]string{}, false, false)
	assert.Len(t, tags, 0)
}

func Test_ParseHTML_NoStopWords(t *testing.T) {
	tags := ParseHTML(html, false, false)
	assert.Len(t, tags, 7)
	assert.Subset(t, ToStrings(tags), []string{"there", "was", "a", "boy", "who's", "name", "jim"})
}

func Test_ParseHTML_WithStopWords(t *testing.T) {
	tags := ParseHTML(html, false, true)
	assert.Len(t, tags, 2)
	assert.Subset(t, ToStrings(tags), []string{"boy", "jim"})
}
