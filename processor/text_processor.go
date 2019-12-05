package processor

import (
	"fmt"
	"strings"
)

// ParseText parses given text lines of text into a slice of tags.
func ParseText(in InputReader, verbose, noStopWords bool) []*Tag {
	if verbose {
		fmt.Println("parsing plain text...")
	}

	lines, err := in.ReadLines()
	if err != nil {
		return []*Tag{}
	}

	if verbose {
		fmt.Printf("got %d lines\n", len(lines))
	}

	if len(lines) == 0 {
		return []*Tag{}
	}

	tokens := make([]string, 0)
	for _, line := range lines {
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
