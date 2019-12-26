package processor

import (
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
	<p>Who's name was Jim.</p>
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
	<p>Who's name was Jim.</p>
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
var parseTests = []struct {
	name    string
	in      string
	expect  []string
	exclude bool
}{
	{
		"empty",
		"",
		[]string{},
		false,
	},
	{
		"simple",
		htmlSimpleString,
		[]string{"there", "was", "a", "boy", "who's", "name", "jim"},
		false,
	},
	{
		"simple exclude stopWords",
		htmlSimpleString,
		[]string{"boy", "jim"},
		true,
	},
	{
		"complex",
		htmlComplexString,
		[]string{"parse", "content", "from", "certain", "tags", "go", "golang", "html", "extract", "all", "theme", "help", "blog"},
		false,
	},
	{
		"complex exclude stopWords",
		htmlComplexString,
		[]string{"parse", "content", "tags", "golang", "html", "extract", "theme", "help", "blog"},
		true,
	},
}

func Test_ParseHTML(t *testing.T) {
	for _, tt := range parseTests {
		t.Run(tt.in, func(t *testing.T) {
			out := ParseHTML(&inputReadCloser{strings.NewReader(tt.in)}, false, tt.exclude)
			assert.ElementsMatch(t, tt.expect, ToStrings(out))
		})
	}
}

func Test_ParseHTML_DedupeTitleAndHeading(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader(htmlDupedString)}, false, true)
	assert.Contains(t, tags, &Tag{Value: "story", Score: 3.0, Count: 1, Docs: 1, DocsCount: 5})
}

func Test_ParseHTML_NoSpecificStopWords(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader(htmlDupedString)}, false, true)
	assert.NotContains(t, tags, &Tag{Value: "part", Score: 1.4, Count: 1})
}
