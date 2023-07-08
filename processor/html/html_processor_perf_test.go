package html

import (
	"bytes"
	"os"
	"testing"

	"github.com/zoomio/tagify/config"
)

var (
	vergeHTML, _   = os.ReadFile("../../_resources_test/html/theverge.html")
	chineseHTML, _ = os.ReadFile("../../_resources_test/html/chinese.html")
)

func BenchmarkParseHTML(b *testing.B) {
	// b.ResetTimer()

	// setup
	cfg, contents := setup(vergeHTML)
	for i := 0; i < b.N; i++ {
		_, _ = tagifyHTML(contents, cfg, nil)
	}
}

func BenchmarkParseHTML_chinese(b *testing.B) {
	cfg, contents := setup(chineseHTML)
	for i := 0; i < b.N; i++ {
		_, _ = tagifyHTML(contents, cfg, nil)
	}
}

func setup(htmlPage []byte) (*config.Config, *HTMLContents) {
	cfg := &config.Config{TagWeights: defaultTagWeights}
	contents := ParseHTML(
		bytes.NewBuffer(htmlPage),
		&config.Config{TagWeights: defaultTagWeights},
		nil,
		nil,
	)
	return cfg, contents
}
