package util

import (
	"math"

	"github.com/zoomio/tagify/processor/model"
)

// TFIDF applies TF-IDF to given Tag
func TFIDF(t *model.Tag) float64 {
	tf := math.Log(1.0 + t.Score)
	idf := math.Log(float64(t.DocsCount) / float64(t.Docs))
	return tf * idf
}
