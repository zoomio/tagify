package processor

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html/atom"
)

const (
	htmlSimpleString = `
	<html>
	<body>
	<p>There was a boy</p>
	<p>Whose name was Jim.</p>
	</body>
	</html>
`

	htmlComplexString = `<!DOCTYPE html>
  <html itemscope itemtype="http://schema.org/QAPage">
  <head>
  	<title>go - Golang parse HTML, extract all content from certain HTML tags</title>
    <link rel="shortcut icon" href="//cdn.sstatic.net/Sites/stackoverflow/img/favicon.ico?v=4f32ecc8f43d">
    <link rel="apple-touch-icon image_src" href="//cdn.sstatic.net/Sites/stackoverflow/img/apple-touch-icon.png?v=c78bd457575a">
    <link rel="search" type="application/opensearchdescription+xml" title="Stack Overflow" href="/opensearch.xml">
    <meta name="twitter:card" content="summary">
    <meta name="twitter:domain" content="stackoverflow.com"/>
    <meta property="og:type" content="website" />
    </head>
<body class="template-blog">
<nav class="navigation">
<div class="navigation__container container">
<a class="navigation__logo" href="/">
<h1>Theme</h1>
</a>
<ul class="navigation__menu">
<li><a href="/help/">Help</a></li>
<li><a href="/blog">Blog</a></li>
</ul>
</div>`

	htmlDupedString = `
	<html>
	<head>
	<title>A story about a boy</title>
	</head>
	<body>
	<h1>A story about a boy</h1>
	<h2>Part I</h2>
	<p>There was a boy</p>
	<p>Whose name was Jim.</p>
	</body>
	</html>
`
)

type inputReadCloser struct {
	io.Reader
}

func (in *inputReadCloser) Close() error {
	return nil
}

// table driven tests
var parseHTMLTests = []struct {
	name        string
	in          string
	expect      []string
	title       string
	hash        string
	noStopWords bool
}{
	{
		"empty",
		"",
		[]string{},
		"",
		"",
		false,
	},
	{
		"simple",
		htmlSimpleString,
		[]string{"there", "was", "a", "boy", "whose", "name", "jim"},
		"",
		"1f4911e9a610990862bbdf6fe1196a4d4003f12896ab0ed20ece0b97fae54bd798ee349bde89e2fd23ccca0063feccd109a4d0d6514f2f0839ff6ac76489bc87",
		false,
	},
	{
		"simple exclude stopWords",
		htmlSimpleString,
		[]string{"boy", "jim"},
		"",
		"1f4911e9a610990862bbdf6fe1196a4d4003f12896ab0ed20ece0b97fae54bd798ee349bde89e2fd23ccca0063feccd109a4d0d6514f2f0839ff6ac76489bc87",
		true,
	},
	{
		"complex",
		htmlComplexString,
		[]string{"parse", "content", "from", "certain", "tags", "go", "golang", "html", "extract", "all"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"e58f7951dca123391bcae296ccbde6abb814a2ef581225caf2f2c6765d39f7da77082a3a1870b6231fcd59e35863dff45e65cd732e1c46911c15269ac021f857",
		false,
	},
	{
		"complex exclude stopWords",
		htmlComplexString,
		[]string{"parse", "content", "tags", "golang", "html", "extract"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"e58f7951dca123391bcae296ccbde6abb814a2ef581225caf2f2c6765d39f7da77082a3a1870b6231fcd59e35863dff45e65cd732e1c46911c15269ac021f857",
		true,
	},
}

func Test_ParseHTML(t *testing.T) {
	for _, tt := range parseHTMLTests {
		t.Run(tt.name, func(t *testing.T) {
			out := ParseHTML(&inputReadCloser{strings.NewReader(tt.in)}, NoStopWords(tt.noStopWords))
			assert.Equal(t, tt.title, out.DocTitle)
			assert.Equal(t, tt.hash, fmt.Sprintf("%x", out.DocHash))
			assert.ElementsMatch(t, tt.expect, ToStrings(out.Tags))
		})
	}
}

func Test_ParseHTML_DedupeTitleAndHeading(t *testing.T) {
	out := ParseHTML(&inputReadCloser{strings.NewReader(htmlDupedString)}, NoStopWords(true))
	assert.Equal(t, "A story about a boy", out.DocTitle)
	assert.Equal(t,
		"4f652c47205d3b922115eef155c484cf81096351696413c86277fa0ed89ebfefe30f81ef6fc6a9d7d654a9292c3cb7aa6f3696052e53c113785a9b1b3be7d4a8",
		fmt.Sprintf("%x", out.DocHash))
	assert.Contains(t, out.Tags, &Tag{Value: "story", Score: 3.0, Count: 1, Docs: 1, DocsCount: 5})
}

func Test_ParseHTML_NoSpecificStopWords(t *testing.T) {
	out := ParseHTML(&inputReadCloser{strings.NewReader(htmlDupedString)}, NoStopWords(true))
	assert.Equal(t, "A story about a boy", out.DocTitle)
	assert.Equal(t,
		"4f652c47205d3b922115eef155c484cf81096351696413c86277fa0ed89ebfefe30f81ef6fc6a9d7d654a9292c3cb7aa6f3696052e53c113785a9b1b3be7d4a8",
		fmt.Sprintf("%x", out.DocHash))
	assert.NotContains(t, out.Tags, &Tag{Value: "part", Score: 1.4, Count: 1})
}

func Test_parseHTML(t *testing.T) {
	const htmlPage = `
	<html>
	<body>
	<p>There was a boy <b>whose</b> name was Jim.</p>
	</body>
	</html>
`
	contents := parseHTML(&inputReadCloser{strings.NewReader(htmlPage)})
	assert.NotNil(t, contents)

	assert.Len(t, contents.lines, 1)

	line := contents.lines[0]
	assert.Len(t, line.parts, 3)

	assert.Equal(t, atom.P, line.parts[0].tag)
	assert.Equal(t, "There was a boy ", string(line.parts[0].data))

	assert.Equal(t, atom.B, line.parts[1].tag)
	assert.Equal(t, "whose", string(line.parts[1].data))

	assert.Equal(t, atom.P, line.parts[2].tag)
	assert.Equal(t, " name was Jim.", string(line.parts[2].data))
}
