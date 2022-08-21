package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var splitTextToWordsTests = []struct {
	name     string
	input    string
	expected string
}{
	{
		"English",
		"Natural language processing includes: tokeniziation, term frequency - inverse term frequency, nearest neighbors, part of speech tagging and many more.",
		"natural| |language| |processing| |includes|:| |tokeniziation|,| |term| |frequency| |-| |inverse| |term| |frequency|,| |nearest| |neighbors|,| |part| |of| |speech| |tagging| |and| |many| |more|.",
	},
	{
		"Chinese",
		"世界有七十亿人口",
		"世界|有|七十|亿|人口",
	},
}

func Test_SplitTextToWords(t *testing.T) {
	for _, tt := range splitTextToWordsTests {
		t.Run(tt.name, func(t *testing.T) {
			text := []byte(tt.input)
			words := SplitTextToWords(text, nil)
			assert.Equal(t, tt.expected, strings.Join(BytesToStrings(words), "|"))
		})
	}
}
