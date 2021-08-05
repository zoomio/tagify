package html

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
)

const (
	htmlWithImg = `
	<html>
	<body>
	<h2>What a beautiful day today!</h2>
	<div><img src="https://example.com/example.png" alt="nature" /></div>
	</body>
	</html>
`
)

func Test_Parse_img(t *testing.T) {
	cfg := config.New(
		config.Verbose(true),
		config.TagWeightsString("h2:1|img:0"),
		config.Extensions([]extension.Extension{newTestImgCrawlerExt()}),
	)
	out := ParseHTML(cfg, &inputReadCloser{strings.NewReader(htmlWithImg)})
	assert.Len(t, out.Extensions, 1)
	res := out.Extensions[0]
	assert.Nil(t, res.Err)
	assert.Equal(t, "test-img-crawler", res.Name)
	assert.Equal(t, "v0.0.1", res.Version)
	images, ok := res.Data["images"]
	assert.True(t, ok)
	srcs, ok := images.([]string)
	assert.True(t, ok)
	assert.Len(t, srcs, 1)
	assert.Equal(t, "https://example.com/example.png", srcs[0])
}

func newTestImgCrawlerExt() *TestImgCrawlerExt {
	return &TestImgCrawlerExt{
		images: []string{},
	}
}

type TestImgCrawlerExt struct {
	images []string
}

func (ext *TestImgCrawlerExt) Name() string {
	return "test-img-crawler"
}

func (ext *TestImgCrawlerExt) Version() string {
	return "v0.0.1"
}

func (ext *TestImgCrawlerExt) Result() *extension.Result {
	return extension.NewResult(ext, map[string]interface{}{"images": ext.images}, nil)
}

func (ext *TestImgCrawlerExt) ParseTag(cfg *config.Config, token *html.Token, lineIdx int) error {
	if token.DataAtom.String() == "img" {
		for _, v := range token.Attr {
			if v.Key == "src" {
				ext.images = append(ext.images, v.Val)
			}
		}
	}
	return nil
}
