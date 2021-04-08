package md

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoomio/tagify/processor/model"
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

type inputReadCloser struct {
	io.Reader
}

func (in *inputReadCloser) Close() error {
	return nil
}

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
		"befd44d2460f1fafe42c37b35f69db984cda479bc51873c9ba94e534c4729a34791822ee76a2aaf40f09d0d4245f3cfd955ddf6f32e146f14fdaf2918ccf2d72",
	},
	{
		"medium",
		mdMediumText,
		true,
		[]string{"boy", "friends", "tea", "ham", "story", "jim", "good", "cakes", "jam", "slices", "delicious"},
		"A story about Jim",
		"f2a797a37071e103a228a53f7c6f040bb6e3615331519b9d63371138d30947010fbb0d67f33fcfac90abcda17891d7e6d241f611c4b9ef0e4721c51a92551976",
	},
	{
		"complex",
		mdComplexText,
		true,
		[]string{"things", "example", "worth", "bar", "stripes", "https", "list", "test", "adding", "dog", "cat", "complex", "text", "link", "foo", "bee"},
		"Complex text for Test",
		"57ab424375a4dbb974e7b2a7af76b71aa10f04db8b8aca9f048519ef0057a581fd3675289d8e567423a37e32854819965e5088d1557cfd72f1c4d5a8c826e732",
	},
}

func Test_ParseMD(t *testing.T) {
	for _, tt := range parseMDTests {
		t.Run(tt.name, func(t *testing.T) {
			out := ParseMD(&inputReadCloser{strings.NewReader(tt.text)}, model.NoStopWords(tt.noStopWords))
			assert.Equal(t, tt.title, out.DocTitle)
			assert.Equal(t, tt.hash, fmt.Sprintf("%x", out.DocHash))
			assert.ElementsMatch(t, tt.tags, model.ToStrings(out.FlatTags()))
		})
	}
}

func Test_mdContents_sentences(t *testing.T) {
	contents := &mdContents{
		lines: []*mdLine{
			{tag: paragraph, data: []byte("There was a boy"), parts: []*mdPart{{tag: paragraph, pos: 0, len: 18}}},
			{tag: paragraph, data: []byte("Whose name was Jim.	"), parts: []*mdPart{{tag: paragraph, pos: 0, len: 21}}},
		},
	}

	l1 := contents.lines[0]
	ss1 := l1.sentences()
	assert.Len(t, ss1, 1)
	assert.Equal(t, "There was a boy", string(ss1[0].data))

	l2 := contents.lines[1]
	ss2 := l2.sentences()
	assert.Len(t, ss2, 2)
	assert.Equal(t, "Whose name was Jim", string(ss2[0].data))
	assert.Equal(t, "", string(ss2[1].data))
}

func Test_mdContents_sentences2(t *testing.T) {
	line := &mdLine{
		tag: paragraph,
		parts: []*mdPart{
			{tag: bold, pos: 0, len: 48},
			{tag: paragraph, pos: 48, len: 20},
		},
		data: []byte("**Sentence number one. And then, number two.*** And finally, three."),
	}

	sents := line.sentences()
	assert.Len(t, sents, 6)
	assert.Equal(t, "", string(sents[5].data))
}
