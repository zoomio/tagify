package processor

import (
	"math"
)

// tfidf applies TF-IDF to given Tag
func tfidf(t *Tag) float64 {
	idf := math.Log(float64(t.DocsCount) / float64(t.Docs))
	return t.Score * idf
}
