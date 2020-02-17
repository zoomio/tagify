package processor

import (
	"crypto/sha512"
	"fmt"
	"strings"
)

// ParseText parses given text lines of text into a slice of tags.
func ParseText(in InputReader, verbose, noStopWords bool) ([]*Tag, []byte) {
	if verbose {
		fmt.Println("parsing plain text...")
	}

	var docsCount int

	lines, err := in.ReadLines()
	if err != nil {
		return []*Tag{}, nil
	}

	if verbose {
		fmt.Printf("got %d lines\n", len(lines))
	}

	if len(lines) == 0 {
		return []*Tag{}, nil
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

	return flatten(tokenIndex), hashTokens(tokens)
}

func hashTokens(ts []string) []byte {
	h := sha512.New()
	for _, t := range ts {
		_, _ = h.Write([]byte(t))
	}
	return h.Sum(nil)
}
