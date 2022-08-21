package util

import (
	"strings"
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
	{"yeah", "yeah", true, false},
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
	{"URL", [][]byte{[]byte("https://www.youtube.com/watch?t=296s&v=HHVQUWnOqEU")}, []string{"youtube"}, false},
	{"URL 2", [][]byte{[]byte("https://zoomio.org/foo/bar?bee=dog")}, []string{"zoomio"}, false},
	{"URL 3", [][]byte{[]byte("https://abc.abc")}, []string{"abc"}, false},
	{"URL 4", [][]byte{[]byte("https://abc.com.au")}, []string{"abc"}, false},
	{"URL 5", [][]byte{[]byte("https://www.abc.com.au")}, []string{"abc"}, false},
	{"URL 6", [][]byte{[]byte("https://my.gov.com.au")}, []string{"my", "gov"}, false},
	{"Quoted Batman", [][]byte{[]byte("'the batman'")}, []string{"batman"}, true},
	{"Quoted Batman w stopwords", [][]byte{[]byte("'the batman'")}, []string{"the", "batman"}, false},
	{"City's", [][]byte{[]byte("city's")}, []string{"city's"}, true},
	{"Yeah", [][]byte{[]byte(", oh yeah lotsa fun and sometimes ")}, []string{"yeah", "lotsa", "fun"}, true},
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

var splitToSentencesTests = []struct {
	name   string
	in     string
	expect string
}{
	{
		"Split",
		"This sentence has a comma, so it'll be split into two halves. This sentence has nothing. Should it though?",
		"This sentence has a comma|so it'll be split into two halves|This sentence has nothing|Should it though",
	},
	{
		"Multiple commas",
		"Natural language processing includes: tokeniziation, term frequency - inverse term frequency, nearest neighbors, part of speech tagging and many more.",
		"Natural language processing includes|tokeniziation|term frequency - inverse term frequency|nearest neighbors|part of speech tagging and many more",
	},
	{
		"Split with word \"yeah\"",
		"Testing is a funny thing, sometimes it is a funtivity, oh yeah lotsa fun and sometimes it is just a drag and an extra hussle. Yup funny, isn't it?",
		"Testing is a funny thing|sometimes it is a funtivity|oh yeah lotsa fun and sometimes it is just a drag and an extra hussle|Yup funny|isn't it",
	},
}

func Test_SplitToSentences(t *testing.T) {
	for _, tt := range splitToSentencesTests {
		t.Run(tt.name, func(t *testing.T) {
			sentences := SplitToSentences([]byte(tt.in))
			assert.Equal(t, tt.expect, strings.Join(BytesToStrings(sentences), "|"))
		})
	}
}
