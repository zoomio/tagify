package html

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
	"github.com/zoomio/tagify/model"
)

const (
	htmlWithImg = `
	<html>
	<body>
	<h2>What a beautiful day today!</h2>
	<div><img src="https://example.com/example.png" alt="nature" /></div>
	<h2>Some ducks in the pond</h2>
	<div><img src="https://example.com/example2.png" alt="nature 2" /></div>
	<h2>Stunning Sunset</h2>
	<div><img src="https://example.com/example3.png" alt="nature 2" /></div>
	</body>
	</html>
`
)

func Test_Ext_Parse_img(t *testing.T) {
	cfg := config.New(
		config.TagWeightsString("h2:1|img:0"),
		config.Extensions([]extension.Extension{newTestImgCrawlerExt()}),
	)
	out := ProcessHTML(cfg, &inputReadCloser{strings.NewReader(htmlWithImg)})
	assert.Len(t, out.Extensions, 1)
	results := out.FindExtResults("test-img-crawler", "v0.0.1") //Extensions["test-img-crawler"]
	assert.Len(t, results, 1)
	res := results[0]
	images, ok := res.Data["images"]
	assert.True(t, ok)
	srcs, ok := images.([]string)
	assert.True(t, ok)
	assert.Len(t, srcs, 3)
	assert.Equal(t, "https://example.com/example.png", srcs[0])
	assert.Equal(t, "https://example.com/example2.png", srcs[1])
	assert.Equal(t, "https://example.com/example3.png", srcs[2])
}

func Test_Ext_Tagify_stopwords(t *testing.T) {
	cfg1 := config.New()
	out1 := ProcessHTML(cfg1, &inputReadCloser{strings.NewReader(htmlWithImg)})

	cfg2 := config.New(config.Extensions([]extension.Extension{&testExtraStopWordsExt{stopWords: []string{"day", "sunset"}}}))
	out2 := ProcessHTML(cfg2, &inputReadCloser{strings.NewReader(htmlWithImg)})

	assert.Len(t, out1.Extensions, 0)
	assert.Len(t, out2.Extensions, 1)

	assert.Contains(t, out1.RawTags, "day")
	assert.Contains(t, out1.RawTags, "sunset")

	assert.NotContains(t, out2.RawTags, "day")
	assert.NotContains(t, out2.RawTags, "sunset")
}

func Test_Ext_ParseEnd(t *testing.T) {
	// Total amount of tags is 7 (see assert.Len)
	cfg := config.New(
		config.ExtraTagWeightsString("img:0"),
		config.NoStopWords(true),
	)
	out := ProcessHTML(cfg, &inputReadCloser{strings.NewReader(htmlWithImg)})
	assert.Equal(t, 7, out.RawLen())

	// Amount of tags when stopped is 3 (see assert.Len)
	cfg = config.New(
		config.ExtraTagWeightsString("img:0"),
		config.Extensions([]extension.Extension{&testStopExt{}}),
		config.NoStopWords(true),
	)
	out = ProcessHTML(cfg, &inputReadCloser{strings.NewReader(htmlWithImg)})
	assert.Equal(t, 3, out.RawLen())
}

func newTestImgCrawlerExt() *testImgCrawlerExt {
	return &testImgCrawlerExt{
		images: []string{},
	}
}

type testImgCrawlerExt struct {
	images []string
}

func (ext *testImgCrawlerExt) Name() string {
	return "test-img-crawler"
}

func (ext *testImgCrawlerExt) Version() string {
	return "v0.0.1"
}

func (ext *testImgCrawlerExt) Result() *extension.Result {
	return extension.NewResult(ext, map[string]interface{}{"images": ext.images}, nil)
}

func (ext *testImgCrawlerExt) ParseTag(cfg *config.Config, token *html.Token, lineIdx int, cnts *HTMLContents) (bool, error) {
	if token.DataAtom.String() == "img" {
		for _, v := range token.Attr {
			if v.Key == "src" {
				ext.images = append(ext.images, v.Val)
			}
		}
	}
	return false, nil
}

type testExtraStopWordsExt struct {
	stopWords []string
}

func (ext *testExtraStopWordsExt) Name() string {
	return "test-extra-stopwords"
}

func (ext *testExtraStopWordsExt) Version() string {
	return "v0.0.1"
}

func (ext *testExtraStopWordsExt) Result() *extension.Result {
	return extension.NewResult(ext, map[string]interface{}{"stopwords": ext.stopWords}, nil)
}

func (ext *testExtraStopWordsExt) Tagify(cfg *config.Config, line *HTMLLine, tokenIndex map[string]*model.Tag) error {
	for _, v := range ext.stopWords {
		delete(tokenIndex, v)
	}
	return nil
}

type testStopExt struct {
}

func (ext *testStopExt) Name() string {
	return "test-stop"
}

func (ext *testStopExt) Version() string {
	return "v0.0.1"
}

func (ext *testStopExt) Result() *extension.Result {
	return extension.NewResult(ext, map[string]interface{}{}, nil)
}

func (ext *testStopExt) ParseTag(cfg *config.Config, token *html.Token, lineIdx int, cnts *HTMLContents) (bool, error) {
	if token.DataAtom.String() == "img" {
		return false, NewHTMLParseEndError()
	}
	return false, nil
}
