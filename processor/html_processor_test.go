package processor

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const htmlSimpleString = `
	<html>
	<body>
	<p>There was a boy</p>
	<p>Who's name was Jim.</p>
	</body>
	</html>
`

type inputReadCloser struct {
	io.Reader
}

func (in *inputReadCloser) Close() error {
	return nil
}

func Test_ParseHTML_Empty(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader("")}, false, false)
	assert.Len(t, tags, 0)
}

func Test_ParseHTML_AllowStopWords(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader(htmlSimpleString)}, false, false)
	assert.Len(t, tags, 7)
	assert.Subset(t, ToStrings(tags), []string{"there", "was", "a", "boy", "who's", "name", "jim"})
}

func Test_ParseHTML_ExcludeStopWords(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader(htmlSimpleString)}, false, true)
	assert.Len(t, tags, 2)
	assert.Subset(t, ToStrings(tags), []string{"boy", "jim"})
}

const htmlComplexString = `<!DOCTYPE html>
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

func Test_ParseHTML_Complex(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader(htmlComplexString)}, false, false)
	assert.Len(t, tags, 11)
	assert.Subset(t, ToStrings(tags), []string{"html", "all", "theme", "parse", "golang", "extract", "content", "from", "certain", "tags", "go"})
}

func Test_ParseHTML_Complex_ExcludeStopWords(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader(htmlComplexString)}, false, true)
	assert.Len(t, tags, 7)
	assert.Subset(t, ToStrings(tags), []string{"golang", "parse", "html", "extract", "content", "tags", "theme"})
}

const htmlDupedString = `
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

func Test_ParseHTML_DedupeTitleAndHeading(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader(htmlDupedString)}, false, true)
	assert.Contains(t, tags, &Tag{Value: "story", Score: 3.0, Count: 1, Docs: 1, DocsCount: 6})
}

func Test_ParseHTML_NoSpecificStopWords(t *testing.T) {
	tags := ParseHTML(&inputReadCloser{strings.NewReader(htmlDupedString)}, false, true)
	fmt.Printf("%v\n", tags)
	assert.NotContains(t, tags, &Tag{Value: "part", Score: 1.4, Count: 1})
}
