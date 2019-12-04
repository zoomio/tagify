package processor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
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
  <title>go - Golang parse HTML, extract all content with &lt;body&gt; &lt;/body&gt; tags - Stack Overflow</title>
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
<h1>Foobar</h1>
</a>
<ul class="navigation__menu">
<li><a href="/tags/">Topics</a></li>
<li><a href="/about">About</a></li>
</ul>
</div>`
)

func Test_ParseHTML_Empty(t *testing.T) {
	tags := ParseHTML(strings.NewReader(""), false, false)
	assert.Len(t, tags, 0)
}

func Test_ParseHTML_AllowStopWords(t *testing.T) {
	tags := ParseHTML(strings.NewReader(htmlSimpleString), false, false)
	assert.Len(t, tags, 7)
	assert.Subset(t, ToStrings(tags), []string{"there", "was", "a", "boy", "who's", "name", "jim"})
}

func Test_ParseHTML_ExcludeStopWords(t *testing.T) {
	tags := ParseHTML(strings.NewReader(htmlSimpleString), false, true)
	assert.Len(t, tags, 2)
	assert.Subset(t, ToStrings(tags), []string{"boy", "jim"})
}

func Test_ParseHTML_Complex(t *testing.T) {
	tags := ParseHTML(strings.NewReader(htmlComplexString), false, false)
	assert.Len(t, tags, 3)
	assert.Subset(t, ToStrings(tags), []string{"topics", "about", "foobar"})
}
