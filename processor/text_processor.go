package processor

import (
	"strings"
)

// ParseText parses given lines of text into a slice of tags.
func ParseText(lines []string, noStopWords bool) []*Tag {
	tokens := make([]string, 0)
	for _, l := range lines {
		tokens = append(tokens, sanitize(strings.Fields(l), noStopWords)...)
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
