package processor

// Tag holds some arbitrary string value (e.g. a word) along with some extra data about it.
type Tag struct {
	Value string
	Score float64
	Count int
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
