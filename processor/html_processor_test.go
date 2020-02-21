package processor

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		[]string{"there", "was", "a", "boy", "who's", "name", "jim"},
		"",
		"d6cdf4700991a6c2db3bb8bb5d2fb57f15e5f0dbe1fcb893781a2d4782b73b43a2232beff70f2ba293599c0d6c8729c4db8a693fdc5dcabb6c10dadca2e31044",
		false,
	},
	{
		"simple exclude stopWords",
		htmlSimpleString,
		[]string{"boy", "jim"},
		"",
		"d6cdf4700991a6c2db3bb8bb5d2fb57f15e5f0dbe1fcb893781a2d4782b73b43a2232beff70f2ba293599c0d6c8729c4db8a693fdc5dcabb6c10dadca2e31044",
		true,
	},
	{
		"complex",
		htmlComplexString,
		[]string{"parse", "content", "from", "certain", "tags", "go", "golang", "html", "extract", "all", "theme", "help", "blog"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"5eef93885dd249586a5f0ae5b03ba02dccfebd18bab9cf0896f891e7b351f62329a13bc5559f8210290f5327d1b5173502437d76eadbed31c3cd7a6e24391958",
		false,
	},
	{
		"complex exclude stopWords",
		htmlComplexString,
		[]string{"parse", "content", "tags", "golang", "html", "extract", "theme", "help", "blog"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"5eef93885dd249586a5f0ae5b03ba02dccfebd18bab9cf0896f891e7b351f62329a13bc5559f8210290f5327d1b5173502437d76eadbed31c3cd7a6e24391958",
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
		"0027df9158090fbd840bf4fe432af56b15ae3d2c460a9b5e2671ed54cbbd8ca75ff803ebbbba7cc2784c18beca10466f3d3a1a954c3f22fcbf66ccc18c751c7b",
		fmt.Sprintf("%x", out.DocHash))
	assert.Contains(t, out.Tags, &Tag{Value: "story", Score: 3.0, Count: 1, Docs: 1, DocsCount: 4})
}

func Test_ParseHTML_NoSpecificStopWords(t *testing.T) {
	out := ParseHTML(&inputReadCloser{strings.NewReader(htmlDupedString)}, NoStopWords(true))
	assert.Equal(t, "A story about a boy", out.DocTitle)
	assert.Equal(t,
		"0027df9158090fbd840bf4fe432af56b15ae3d2c460a9b5e2671ed54cbbd8ca75ff803ebbbba7cc2784c18beca10466f3d3a1a954c3f22fcbf66ccc18c751c7b",
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
	contents := parseHTML(&inputReadCloser{strings.NewReader(htmlDupedString)})

}
