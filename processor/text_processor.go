package processor

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"strings"
)

// ParseText parses given text lines of text into a slice of tags.
var ParseText ParseFunc = func(in io.ReadCloser, options ...ParseOption) *ParseOutput {

	c := &parseConfig{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	if c.verbose {
		fmt.Println("parsing plain text...")
	}

	var docsCount int

	defer in.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)
	inStr := buf.String()
	lines := strings.FieldsFunc(inStr, func(r rune) bool {
		return r == '\n'
	})

	if c.verbose {
		fmt.Printf("got %d lines\n", len(lines))
	}

	if len(lines) == 0 {
		return &ParseOutput{}
	}

	tokenIndex := make(map[string]*Tag)
	tokens := make([]string, 0)
	for _, l := range lines {
		sentences := SplitToSentences([]byte(l))
		for _, s := range sentences {
			docsCount++
			tokens = append(tokens, sanitize(bytes.Fields(s), c.noStopWords)...)
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

	return &ParseOutput{Tags: flatten(tokenIndex), DocHash: hashTokens(tokens)}
}

func hashTokens(ts []string) []byte {
	h := sha512.New()
	for _, t := range ts {
		_, _ = h.Write([]byte(t))
	}
	return h.Sum(nil)
}
