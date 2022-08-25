package config

import (
	"bytes"

	"github.com/go-ego/gse"
)

type Segmenter interface {
	Segment(text []byte) [][]byte
}

type DefaultSegmenter struct {
	lang string
	seg  gse.Segmenter
}

func NewDefaultSegmenter(c *Config) *DefaultSegmenter {
	seg := gse.Segmenter{SkipLog: true}
	if c != nil {
		seg.NotStop = c.NoStopWords
		seg.SkipLog = !c.Verbose
		if c.StopWords != nil {
			seg.StopWordMap = c.StopWords.Index()
		}
	}
	seg.LoadDict()
	s := &DefaultSegmenter{seg: seg}
	if c != nil {
		s.lang = c.Lang
	}
	return s
}

func (s *DefaultSegmenter) Segment(text []byte) [][]byte {
	if s.lang == "zh" || s.lang == "ja" {
		segments := s.seg.Segment(text)
		bs := make([][]byte, len(segments))
		for k, v := range segments {
			bs[k] = []byte(v.Token().Text())
		}
		return bs
	}
	return bytes.Fields(text)
}

func BytesToStrings(txts [][]byte) []string {
	strs := make([]string, len(txts))
	for i := 0; i < len(txts); i++ {
		strs[i] = string(txts[i])
	}
	return strs
}
