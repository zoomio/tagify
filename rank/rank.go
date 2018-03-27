package rank

// Item holds some arbitrary string value (e.g. a word) along with some extra data about it.
type Item struct {
	Value string
	Score float64
}

func flatten(dict map[string]*Item) []*Item {
	flat := make([]*Item, len(dict))
	var i int
	for _, val := range dict {
		flat[i] = val
		i++
	}
	return flat
}
