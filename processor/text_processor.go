package processor

import (
	"strings"
)

// ParseText parses given text lines of text into a slice of tags.
func ParseText(text []string, noStopWords bool) []*Tag {
	tokens := make([]string, 0)
	for _, line := range text {
		tokens = append(tokens, sanitize(strings.Fields(line), noStopWords)...)
	}

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
