package util

import (
	"bytes"

	"github.com/go-ego/gse"
	"github.com/zoomio/stopwords"

	"github.com/zoomio/tagify/config"
)

func SplitTextToWords(text []byte, cfg *config.Config) [][]byte {
	seg := gse.Segmenter{SkipLog: true}
	if cfg != nil {
		seg.NotStop = cfg.NoStopWords
		seg.SkipLog = !cfg.Verbose
		if cfg.StopWords != nil {
			seg.StopWordMap = cfg.StopWords.Index()
		}
	}
	seg.LoadDict()

	segments := seg.Segment(text)

	bs := make([][]byte, len(segments))
	for k, v := range segments {
		bs[k] = []byte(v.Token().Text())
	}

	return bs
}

func BytesToStrings(txts [][]byte) []string {
	strs := make([]string, len(txts))
	for i := 0; i < len(txts); i++ {
		strs[i] = string(txts[i])
	}
	return strs
}

func SplitToTokens(text []byte, cfg *config.Config, lang string, reg *stopwords.Register) []string {
	if lang == "zh" || lang == "ja" {
		return Sanitize(SplitTextToWords(text, cfg), reg)
	}
	return Sanitize(bytes.Fields(text), reg)
}
