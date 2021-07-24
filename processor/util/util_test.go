package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoomio/stopwords"
)

var register = stopwords.Setup()

// table driven tests
var normalizeTests = []struct {
	in      string
	expect  string
	pass    bool
	exclude bool
}{
	{"(part)", "part", true, false},
	{"(part1)", "part", true, false},
	{"part1", "part", true, false},
	{"part}1", "part", true, false},
	{"{part}", "part", true, false},
	{"{d'arko}", "d'arko", true, false},
	{"1)", "1)", false, false},
	{"-no-stop", "-no-stop", false, false},
	{"2018-02-24T12:00:49Z", "--TZ", false, false},
}

func Test_normalize(t *testing.T) {
	for _, tt := range normalizeTests {
		t.Run(tt.in, func(t *testing.T) {
			var reg *stopwords.Register
			if tt.exclude {
				reg = register
			}
			out, ok := Normalize(tt.in, reg)
			assert.Equal(t, tt.pass, ok)
			assert.Equal(t, tt.expect, out)
		})
	}
}

// table driven tests
var sanitizeTests = []struct {
	name    string
	in      [][]byte
	expect  []string
	exclude bool
}{
	{"splits", [][]byte{[]byte("Advertising?Programmes")}, []string{"advertising", "programmes"}, false},
	{"apostrophe", [][]byte{[]byte("I’ve")}, []string{}, true},
}

func Test_sanitize(t *testing.T) {
	for _, tt := range sanitizeTests {
		t.Run(tt.name, func(t *testing.T) {
			var reg *stopwords.Register
			if tt.exclude {
				reg = register
			}
			out := Sanitize(tt.in, reg)
			assert.ElementsMatch(t, tt.expect, out)
		})
	}
}

func Test_SplitToSentences(t *testing.T) {
	text := "This sentence has a comma, so it'll be split into two halves. This sentence has nothing. Should it though?"
	sentences := SplitToSentences([]byte(text))
	assert.Len(t, sentences, 4)
}

func Test_SplitToSentences_MultipleCommas(t *testing.T) {
	text := `
	Natural language processing includes: tokeniziation, term frequency - inverse term frequency, nearest neighbors, part of speech tagging and many more.
	`
	sentences := SplitToSentences([]byte(text))
	assert.Len(t, sentences, 6)
	assert.Equal(t, " part of speech tagging and many more", string(sentences[4]))
}
