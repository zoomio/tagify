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

	var docsCount int

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

	tokenIndex := make(map[string]*Tag)
	tokens := make([]string, 0)
	for _, l := range lines {
		sentences := SplitToSentences(l)
		for _, s := range sentences {
			docsCount++
			tokens = append(tokens, sanitize(strings.Fields(s), noStopWords)...)
			visited := map[string]bool{}
			for _, token := range tokens {
				visited[token] = true
				item, ok := tokenIndex[token]
				if !ok {
					item = &Tag{Value: token}
					tokenIndex[token] = item
				}
				item.Score++
				item.Count++
			}
			// increment number of appearances in documents for each visited tag
			for token := range visited {
				tokenIndex[token].Docs++
			}
		}
	}

	// set total number of dicuments in the text.
	for _, v := range tokenIndex {
		v.DocsCount = docsCount
	}

	return flatten(tokenIndex)
}
