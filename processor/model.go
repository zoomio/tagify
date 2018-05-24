package processor

import "fmt"

// Tag holds some arbitrary string value (e.g. a word) along with some extra data about it.
type Tag struct {
	Value string
	Score float64
	Count int
}

func (t *Tag) String() string {
	return fmt.Sprintf("[%.2f, %d] %s", t.Score, t.Count, t.Value)
}

func flatten(dict map[string]*Tag) []*Tag {
	flat := make([]*Tag, len(dict))
	var i int
	for _, val := range dict {
		flat[i] = val
		i++
	}
	return flat
}
