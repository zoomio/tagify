package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	{"2018-02-24T12:00:49Z", "--", false, false},
}

func Test_normalize(t *testing.T) {
	for _, tt := range normalizeTests {
		t.Run(tt.in, func(t *testing.T) {
			out, ok := normalize(tt.in, tt.exclude)
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
			out := sanitize(tt.in, tt.exclude)
			assert.ElementsMatch(t, tt.expect, out)
		})
	}
}

func Test_Run_Limits(t *testing.T) {
	items := []*Tag{
		{Value: "cat", Score: 1},
		{Value: "dog", Score: 1},
		{Value: "foo", Score: 1},
		{Value: "bar", Score: 1},
		{Value: "bee", Score: 1},
	}
	processed := Run(items, 3)
	assert.Len(t, processed, 3)
}

func Test_Run_Sorts(t *testing.T) {
	items := []*Tag{
		{Value: "cat", Score: 1},
		{Value: "dog", Score: 2},
		{Value: "foo", Score: 5},
		{Value: "bar", Score: 3},
		{Value: "bee", Score: 4},
	}
	processed := Run(items, 5)
	assert.Len(t, processed, 5)
	assert.Equal(t, "foo", processed[0].Value)
	assert.Equal(t, "bee", processed[1].Value)
	assert.Equal(t, "bar", processed[2].Value)
	assert.Equal(t, "dog", processed[3].Value)
	assert.Equal(t, "cat", processed[4].Value)
}

func Test_Run_DeDupes(t *testing.T) {
	items := []*Tag{
		{Value: "cat", Score: 5},
		{Value: "person", Score: 2},
		{Value: "people", Score: 5},
		{Value: "bar", Score: 3},
		{Value: "cats", Score: 1},
	}
	processed := Run(items, 5)
	assert.Len(t, processed, 3)
	assert.Equal(t, "people", processed[0].Value)
	assert.Equal(t, "cat", processed[1].Value)
	assert.Equal(t, "bar", processed[2].Value)
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

func Test_Run_IgnoresTFIDF_IfNoDocs(t *testing.T) {
	items := []*Tag{
		{Value: "cat", Score: 5},
	}
	processed := Run(items, 5)
	assert.Equal(t, 5.0, processed[0].Score)
}

func Test_Run_AppliesTFIDF_WithDocs(t *testing.T) {
	items := []*Tag{
		{Value: "cat", Score: 5, Docs: 1, DocsCount: 3},
	}
	processed := Run(items, 5)
	assert.Equal(t, 1.9684489712313906, processed[0].Score)
}
