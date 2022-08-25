package util

import (
	"github.com/zoomio/stopwords"

	"github.com/zoomio/tagify/config"
)

func SplitToTokens(text []byte, cfg *config.Config) []string {
	var reg *stopwords.Register
	if cfg.NoStopWords {
		reg = cfg.StopWords
	}
	return Sanitize(cfg.Segment(text), reg)
}
