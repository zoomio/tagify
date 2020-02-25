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

	cssyHTML = `
<html>
	<head>
		<meta charset="utf-8">
		<meta content="IE=edge" http-equiv="X-UA-Compatible">
		<meta content="width=device-width, initial-scale=1.0, viewport-fit=cover" name="viewport">
		<script async="" src="https://www.googletagmanager.com/gtm.js?id=GTM-M2PTNW9"></script>
		<script defer="" src="/polyfill.min.js?flags=gated&amp;features=IntersectionObserver%2CIntersectionObserverEntry%2Ces2017%2Ces2016"></script>
		<meta name="generator" content="Gatsby 2.19.10">
		<title>People are looking for cool stuff</title>
		<link data-react-helmet="true" rel="shortcut icon" href="/5f1ac6743fee20e0237f57e2a2cf5be5.png" type="image/x-icon">
		<meta data-react-helmet="true" name="twitter:card" content="summary_large_image">
		<meta data-react-helmet="true" name="twitter:site" content="@hifoo">
		<meta data-react-helmet="true" name="twitter:creator" content="@hifoo">
		<meta data-react-helmet="true" name="description" content="Foo.">
		<meta data-react-helmet="true" property="og:type" content="webpage">
		<meta data-react-helmet="true" property="og:image" content="https://foo.com/39b8e6a9a8367644e4b1fbfd74285549.png">
		<meta data-react-helmet="true" property="og:description" content="Foo bar.">
		<meta data-react-helmet="true" property="og:title" content="Bar foo">
		<meta data-react-helmet="true" name="viewport" content="width=device-width, initial-scale=1.0">
		<meta data-react-helmet="true" name="keywords" content="Stuff is hard.">
		<script data-react-helmet="true">if (window.localStorage.getItem("current-user") !== null) {
              document.documentElement.className += " logged-in";
			}</script>
		<script data-react-helmet="true" type="application/ld+json">{"@context":"https://schema.org","@type":"WebPage","description":"Machines some day are gonna rise"}</script>
		<link rel="sitemap" type="application/xml" href="/sitemap.xml">
		<style type="text/css">.gatsby-resp-image-image{width:100%;height:100%;margin:0;vertical-align:middle;position:absolute;top:0;left:0;color:transparent;}</style>
		<script>
			document.addEventListener("DOMContentLoaded", function(event) {
			var hash = window.decodeURI(location.hash.replace('#', ''))
			if (hash !== '') {
				var element = document.getElementById(hash)
				if (element) {
				let scrollTop = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop
				let clientTop = document.documentElement.clientTop || document.body.clientTop || 0
				var offset = element.getBoundingClientRect().top + scrollTop - clientTop
				// Wait for the browser to finish rendering before scrolling.
				setTimeout((function() {
					window.scrollTo(0, offset - 96)
				}), 0)
				}
			}
			})
  		</script>
		  <link rel="alternate" type="application/rss+xml" title="Foo blog" href="/rss.xml">
		  <script>(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start': new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0], j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src= 'https://www.googletagmanager.com/gtm.js?id='+i+dl+'';f.parentNode.insertBefore(j,f); })(window,document,'script','dataLayer', 'GTM-M2PTNW9');</script>
		  <link rel="canonical" href="https://fooapp.com/" data-baseprotocol="https:" data-basehost="fooapp.com">
		  <link as="script" rel="preload" href="/app-898efce57969bc7915fa.js">
		  <link as="script" rel="preload" href="/component---src-pages-index-tsx-a6d43c978e1949ea74f0.js">
		  <link as="script" rel="preload" href="/commons-095679602a6e1870f198.js">
		  <link as="script" rel="preload" href="/webpack-runtime-d4cc4af9dba796da0984.js">
		  <link as="fetch" rel="preload" href="/page-data/index/page-data.json" crossorigin="anonymous">
		  <style data-emotion="css-global"></style>
		  <style data-emotion-css="1vv127u">.css-1vv127u{-webkit-align-items:center;-webkit-box-align:center;-ms-flex-align:center;align-items:center;background-color:#24124d;color:#ffffff;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;font-size:16px;font-weight:500;-webkit-box-pack:center;-webkit-justify-content:center;-ms-flex-pack:center;justify-content:center;line-height:24px;padding:16px;text-align:center;}@media print{.css-1vv127u{display:none;}}@media (max-width:440px){.css-1vv127u{display:none;}}</style><style data-emotion-css="1q9d2m6">.css-1q9d2m6{margin:0 16px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;}</style><style data-emotion-css="11f3jc2">.css-11f3jc2{color:#b9cdfc;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;}.css-11f3jc2:hover,.css-11f3jc2:focus{color:#dce6fe;}</style><style data-emotion-css="pittvc">.css-pittvc{-webkit-align-items:center;-webkit-box-align:center;-ms-flex-align:center;align-items:center;background-color:#ffffff;color:#24124d;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;height:64px;padding:0 24px;position:relative;-webkit-transition:box-shadow 300ms cubic-bezier(.2,.6,.6,1);transition:box-shadow 300ms cubic-bezier(.2,.6,.6,1);z-index:1000;}@media print{.css-pittvc{display:none;}}</style>
		  <style data-emotion-css="uhbzfs">.css-uhbzfs{-webkit-align-items:center;-webkit-box-align:center;-ms-flex-align:center;align-items:center;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-flex:0 1 180px;-ms-flex:0 1 180px;flex:0 1 180px;justify-self:flex-start;line-height:0;}</style><style data-emotion-css="5wf59g">.css-5wf59g{border-radius:4px;display:inline-block;-webkit-transition:box-shadow 150ms cubic-bezier(.2,.6,.6,1);transition:box-shadow 150ms cubic-bezier(.2,.6,.6,1);}.css-5wf59g:hover{background-color:#f8f7fc;}.css-5wf59g:focus{box-shadow:0 0 0 3px rgba(91,147,255,.4);}</style>
		  <style data-emotion-css="15zv5un">.css-15zv5un{border-radius:4px;margin-left:8px;-webkit-transition:box-shadow 150ms cubic-bezier(.2,.6,.6,1);transition:box-shadow 150ms cubic-bezier(.2,.6,.6,1);}.css-15zv5un:hover,.css-15zv5un:focus{opacity:0.9;}.css-15zv5un:focus{box-shadow:0 0 0 3px rgba(91,147,255,.4);}@media (max-width:440px){.css-15zv5un{display:none;}}</style>
		  <style data-emotion-css="1w7oq1a">.css-1w7oq1a{-webkit-align-items:center;-webkit-box-align:center;-ms-flex-align:center;align-items:center;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-flex:1 1 auto;-ms-flex:1 1 auto;flex:1 1 auto;-webkit-box-pack:center;-webkit-justify-content:center;-ms-flex-pack:center;justify-content:center;list-style:none;min-width:0;margin:0 auto;padding:0;}@media (max-width:960px){.css-1w7oq1a{display:none;}}</style>
		  <style data-emotion-css="1ugrdw8">.css-1ugrdw8{-webkit-align-items:center;-webkit-box-align:center;-ms-flex-align:center;align-items:center;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;margin:0;padding:0;}.css-1ugrdw8:hover > a{background-color:rgba(81,45,168,0.04);color:#512da8;}.css-1ugrdw8:focus > a{box-shadow:0 0 0 3px rgba(91,147,255,.4);color:#512da8;}.css-1ugrdw8.active > a{background-color:rgba(81,45,168,0.08);color:#512da8;}.css-1ugrdw8:hover > .nav-submenu,.css-1ugrdw8 > .nav-submenu:hover,.css-1ugrdw8:focus > .nav-submenu,.css-1ugrdw8:focus-within > div{visibility:visible;}.css-1ugrdw8:hover > .nav-submenu > div,.css-1ugrdw8 > .nav-submenu:hover > div,.css-1ugrdw8:focus > .nav-submenu > div,.css-1ugrdw8:focus-within > div > div{opacity:1;-webkit-transform:translate3d(0px,0px,0px);-ms-transform:translate3d(0px,0px,0px);transform:translate3d(0px,0px,0px);}</style>
		  <style data-emotion-css="4ht8ou">.css-4ht8ou{background:transparent;border-radius:4px;border:0;color:inherit;cursor:pointer;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;font-size:16px;font-weight:700;line-height:24px;margin:0 2px;outline:0;overflow:hidden;padding:4px 12px;text-overflow:ellipsis;-webkit-transition:background-color 150ms cubic-bezier(.2,.6,.6,1),box-shadow 150ms cubic-bezier(.2,.6,.6,1),color 150ms cubic-bezier(.2,.6,.6,1);transition:background-color 150ms cubic-bezier(.2,.6,.6,1),box-shadow 150ms cubic-bezier(.2,.6,.6,1),color 150ms cubic-bezier(.2,.6,.6,1);white-space:nowrap;-webkit-appearance:none;z-index:10000;}.css-4ht8ou:hover{background-color:rgba(81,45,168,0.04);color:#512da8;}.css-4ht8ou:focus{box-shadow:0 0 0 3px rgba(91,147,255,.4);color:#512da8;}.css-4ht8ou.active{background-color:rgba(81,45,168,0.08);color:#512da8;}</style>
		  <style data-emotion-css="r7p3o9">.css-r7p3o9{background-color:#ffffff;box-shadow:0 32px 48px -8px rgba(36,18,77,.2);left:0;min-height:380px;padding:32px 16px;position:absolute;right:0;top:calc(100% - 16px);visibility:hidden;width:100%;will-change:visibility;z-index:1000;}.css-r7p3o9 > div{opacity:0;-webkit-transition:opacity 300ms cubic-bezier(.2,.6,.6,1),-webkit-transform 300ms cubic-bezier(.2,.6,.6,1);-webkit-transition:opacity 300ms cubic-bezier(.2,.6,.6,1),transform 300ms cubic-bezier(.2,.6,.6,1);transition:opacity 300ms cubic-bezier(.2,.6,.6,1),transform 300ms cubic-bezier(.2,.6,.6,1);-webkit-transform:translate3d(0px,-16px,0px);-ms-transform:translate3d(0px,-16px,0px);transform:translate3d(0px,-16px,0px);-webkit-transform-style:preserve-3d;-ms-transform-style:preserve-3d;transform-style:preserve-3d;will-change:opacity,transform;}</style>
		  <style data-emotion-css="1pv4qoq">.css-1pv4qoq{display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-flex-direction:column;-ms-flex-direction:column;flex-direction:column;margin:0 auto;max-width:1080px;}</style><style data-emotion-css="k008qs">.css-k008qs{display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;}</style><style data-emotion-css="11qjisw">.css-11qjisw{-webkit-flex:1 1 auto;-ms-flex:1 1 auto;flex:1 1 auto;}</style>
		  <style data-emotion-css="16wtyzb">.css-16wtyzb{box-shadow:0 1px 0 0 #eeecf1;font-size:14px;font-weight:500;-webkit-letter-spacing:1px;-moz-letter-spacing:1px;-ms-letter-spacing:1px;letter-spacing:1px;line-height:24px;margin:0 16px 16px;overflow:hidden;padding:0 0 16px;text-overflow:ellipsis;text-transform:uppercase;white-space:nowrap;}</style><style data-emotion-css="137k06d">.css-137k06d{display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;list-style:none;margin:0;padding:0;}</style>
		  <style data-emotion-css="dagzzh">.css-dagzzh{background:transparent;border:0;border-radius:8px;cursor:pointer;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-flex-direction:column;-ms-flex-direction:column;flex-direction:column;-webkit-box-pack:justify;-webkit-justify-content:space-between;-ms-flex-pack:justify;justify-content:space-between;margin:0;min-height:200px;outline:0;padding:16px;text-align:left;-webkit-transition:background-color 150ms cubic-bezier(.2,.6,.6,1),box-shadow 150ms cubic-bezier(.2,.6,.6,1);transition:background-color 150ms cubic-bezier(.2,.6,.6,1),box-shadow 150ms cubic-bezier(.2,.6,.6,1);width:100%;}.css-dagzzh strong{font-size:18px;font-weight:700;line-height:24px;}.css-dagzzh p{font-size:14px;font-weight:400;line-height:20px;margin-top:4px;}.css-dagzzh:hover{background:#fff9eb;}.css-dagzzh:focus{background:#fff9eb;box-shadow:0 0 0 3px rgba(91,147,255,.4);}.css-dagzzh.active{background:#fff9eb;}</style>
		  <style data-emotion-css="tkst1s">.css-tkst1s{background:transparent;border:0;border-radius:8px;cursor:pointer;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-flex-direction:column;-ms-flex-direction:column;flex-direction:column;-webkit-box-pack:justify;-webkit-justify-content:space-between;-ms-flex-pack:justify;justify-content:space-between;margin:0;min-height:200px;outline:0;padding:16px;text-align:left;-webkit-transition:background-color 150ms cubic-bezier(.2,.6,.6,1),box-shadow 150ms cubic-bezier(.2,.6,.6,1);transition:background-color 150ms cubic-bezier(.2,.6,.6,1),box-shadow 150ms cubic-bezier(.2,.6,.6,1);width:100%;}.css-tkst1s strong{font-size:18px;font-weight:700;line-height:24px;}.css-tkst1s p{font-size:14px;font-weight:400;line-height:20px;margin-top:4px;}.css-tkst1s:hover{background:#f1f5ff;}.css-tkst1s:focus{background:#f1f5ff;box-shadow:0 0 0 3px rgba(91,147,255,.4);}.css-tkst1s.active{background:#f1f5ff;}</style>
	</head>
	<body>
			<div id="content">
				<p>Texty text about very important stuff.</p>
			</div>
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
	contentOnly bool
}{
	{
		"empty",
		"",
		[]string{},
		"",
		"",
		false,
		true,
	},
	{
		"simple",
		htmlSimpleString,
		[]string{"there", "was", "a", "boy", "whose", "name", "jim"},
		"",
		"1f4911e9a610990862bbdf6fe1196a4d4003f12896ab0ed20ece0b97fae54bd798ee349bde89e2fd23ccca0063feccd109a4d0d6514f2f0839ff6ac76489bc87",
		false,
		true,
	},
	{
		"simple exclude stopWords",
		htmlSimpleString,
		[]string{"boy", "jim"},
		"",
		"1f4911e9a610990862bbdf6fe1196a4d4003f12896ab0ed20ece0b97fae54bd798ee349bde89e2fd23ccca0063feccd109a4d0d6514f2f0839ff6ac76489bc87",
		true,
		true,
	},
	{
		"complex",
		htmlComplexString,
		[]string{"parse", "content", "from", "certain", "tags", "go", "golang", "html", "extract", "all"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"0b1e1436f1918ec3e331c9d865d88d8fdd82051dac258658e67a270b0d53b45572fa11a24df322dd2ea4dde8e374b9d33d8ac68940ef8979a13f0ca71d385a4f",
		false,
		true,
	},
	{
		"complex exclude stopWords",
		htmlComplexString,
		[]string{"parse", "content", "tags", "golang", "html", "extract"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"0b1e1436f1918ec3e331c9d865d88d8fdd82051dac258658e67a270b0d53b45572fa11a24df322dd2ea4dde8e374b9d33d8ac68940ef8979a13f0ca71d385a4f",
		true,
		true,
	},
	{
		"complex exclude stopWords tag everything",
		htmlComplexString,
		[]string{"html", "content", "tags", "theme", "blog", "parse", "extract", "help", "golang"},
		"go - Golang parse HTML, extract all content from certain HTML tags",
		"0b1e1436f1918ec3e331c9d865d88d8fdd82051dac258658e67a270b0d53b45572fa11a24df322dd2ea4dde8e374b9d33d8ac68940ef8979a13f0ca71d385a4f",
		true,
		false,
	},
	{
		"css-y",
		cssyHTML,
		[]string{"stuff", "texty", "text", "people", "cool"},
		"People are looking for cool stuff",
		"01443073300b7a758c6cdcd826e04c66b008acf033ef953231bee8119e1f0400e85e4702ce2cd929873e44bb2b0d550fea27bce2014ef24e6c68159e2a170210",
		true,
		false,
	},
}

func Test_ParseHTML(t *testing.T) {
	for _, tt := range parseHTMLTests {
		t.Run(tt.name, func(t *testing.T) {
			out := ParseHTML(&inputReadCloser{strings.NewReader(tt.in)}, NoStopWords(tt.noStopWords), ContentOnly(tt.contentOnly))
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
	assert.Contains(t, out.Tags, &Tag{Value: "story", Score: 3.0, Count: 1, Docs: 1, DocsCount: 4})
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
	contents := parseHTML(&inputReadCloser{strings.NewReader(htmlPage)}, nil)
	assert.NotNil(t, contents)

	assert.Len(t, contents.lines, 1)

	line := contents.lines[0]
	assert.Len(t, line.parts, 3)

	assert.Equal(t, atom.P, line.parts[0].tag)
	assert.Equal(t, "There was a boy ", string(line.pData(line.parts[0])))

	assert.Equal(t, atom.B, line.parts[1].tag)
	assert.Equal(t, "whose", string(line.pData(line.parts[1])))

	assert.Equal(t, atom.P, line.parts[2].tag)
	assert.Equal(t, " name was Jim.", string(line.pData(line.parts[2])))
}
