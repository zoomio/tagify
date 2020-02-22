package processor

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mdSmallText = `
There was a boy
Whose name was Jim.
`

	mdMediumText = `
# A story about Jim
	
There was a boy whose name was **Jim**.	His Friends were very good to him.
They gave him Tea, and Cakes, and Jam, And slices of delicious Ham...
`

	mdComplexText = `
# Complex text for Test
	
This text will contain a **few** things. It will have a [link to the example](https://example.com).

## Very important section

Few words here and __there__. Maybe also worth adding **a list**, right?
Said list:
- Foo, bar
- Bee and dog - or a cat with **stripes**
`
)

// table driven tests
var parseMDTests = []struct {
	name        string
	text        string
	noStopWords bool
	tags        []string
	title       string
	hash        string
}{
	{
		"small",
		mdSmallText,
		false,
		[]string{"was", "there", "a", "boy", "whose", "name", "jim"},
		"",
		"3604d570a9face3d21333b9e15818fb24cb3d5b142b18ad6cb41164798638e758607a06a042739c74873850b8043f8715759740b8a1f2c886ccf9d85d0f159c0",
	},
	{
		"medium",
		mdMediumText,
		true,
		[]string{"boy", "friends", "tea", "ham", "story", "jim", "good", "cakes", "jam", "slices", "delicious"},
		"A story about Jim",
		"ff73e809ba68765670d32ddbb3b1dad8a75bfee83bd30dfce311a16eda9069f08a07b118d17086a3830b6db3bb1be36f059dffada6bb1c2e9eb0e24c34f2d220",
	},
	{
		"complex",
		mdComplexText,
		true,
		[]string{"things", "example", "worth", "bar", "stripes", "https", "list", "test", "adding", "dog", "cat", "complex", "text", "link", "foo", "bee"},
		"Complex text for Test",
		"c1d5d1b299313bc7a44102cd92f4160edc734a32d580890a1a6b9682ac5ee8bbba3d9c226db454b748ec9923045553a115c22eb9be378da5d42c68ef365122c6",
	},
}

func Test_ParseMD(t *testing.T) {
	for _, tt := range parseMDTests {
		t.Run(tt.name, func(t *testing.T) {
			out := ParseMD(&inputReadCloser{strings.NewReader(tt.text)}, NoStopWords(tt.noStopWords))
			assert.Equal(t, tt.title, out.DocTitle)
			assert.Equal(t, tt.hash, fmt.Sprintf("%x", out.DocHash))
			assert.ElementsMatch(t, tt.tags, ToStrings(out.Tags))
		})
	}
}

func Test_mdContents_sentences(t *testing.T) {
	contents := &mdContents{
		lines: []*mdLine{
			{tag: paragraph, parts: []*mdPart{{tag: paragraph, data: []byte("There was a boy")}}},
			{tag: paragraph, parts: []*mdPart{{tag: paragraph, data: []byte("Whose name was Jim.	")}}},
		},
	}

	l1 := contents.lines[0]
	ss1 := l1.sentences()
	assert.Len(t, ss1, 1)
	assert.Equal(t, "There was a boy", string(ss1[0].data()))

	l2 := contents.lines[1]
	ss2 := l2.sentences()
	assert.Len(t, ss2, 2)
	assert.Equal(t, "Whose name was Jim", string(ss2[0].data()))
	assert.Equal(t, "", string(ss2[1].data()))
}

func Test_mdContents_sentences2(t *testing.T) {
	line := &mdLine{
		tag: paragraph,
		parts: []*mdPart{
			{tag: bold, data: []byte("**Sentence number one. And then, number two.***")},
			{tag: paragraph, data: []byte(" And finally, three.")},
		},
	}

	sents := line.sentences()
	assert.Len(t, sents, 6)
	assert.Equal(t, "", string(sents[5].data()))
}
