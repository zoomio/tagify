package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var splitTextToWordsTests = []struct {
	name     string
	lang     string
	input    string
	expected string
}{
	{
		"English",
		"en",
		"Natural language processing includes: tokeniziation, term frequency - inverse term frequency, nearest neighbors, part of speech tagging and many more.",
		"Natural|language|processing|includes:|tokeniziation,|term|frequency|-|inverse|term|frequency,|nearest|neighbors,|part|of|speech|tagging|and|many|more.",
	},
	{
		"Chinese",
		"zh",
		"世界有七十亿人口",
		"世界|有|七十|亿|人口",
	},
}

func Test_SplitTextToWords(t *testing.T) {

	for _, tt := range splitTextToWordsTests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := New(Language(tt.lang))
			text := []byte(tt.input)
			words := cfg.Segment(text)
			assert.Equal(t, tt.expected, strings.Join(BytesToStrings(words), "|"))
		})
	}
}
