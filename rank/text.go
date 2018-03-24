package rank

// ParseText ...
func ParseText(tokens []string) []*Item {
	tokens = Filter(tokens)
	if len(tokens) == 0 {
		return []*Item{}
	}
	index := make(map[string]*Item)
	for _, token := range tokens {
		item, ok := index[token]
		if !ok {
			item = &Item{Value: token}
			index[token] = item
		}
		item.Score = item.Score + 1
	}
	return flatten(index)
}
