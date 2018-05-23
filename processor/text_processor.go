package processor

// ParseText ...
func ParseText(tokens []string) []*Tag {
	tokens = Filter(tokens)
	if len(tokens) == 0 {
		return []*Tag{}
	}
	index := make(map[string]*Tag)
	for _, token := range tokens {
		item, ok := index[token]
		if !ok {
			item = &Tag{Value: token}
			index[token] = item
		}
		item.Score++
		item.Count++
	}
	return flatten(index)
}
