package processor

import (
	"math"
)

// tfidf applies TF-IDF to given Tag
func tfidf(t *Tag) float64 {
	tf := math.Log(1.0 + t.Score)
	idf := math.Log(float64(t.DocsCount) / float64(t.Docs))
	return tf * idf
}
