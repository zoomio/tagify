package text

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"strings"

	"github.com/zoomio/tagify/processor/model"
	"github.com/zoomio/tagify/processor/util"
)

// ParseText parses given text lines of text into a slice of tags.
var ParseText model.ParseFunc = func(in io.ReadCloser, options ...model.ParseOption) *model.ParseOutput {

	c := &model.ParseConfig{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	if c.Verbose {
		fmt.Println("parsing plain text...")
	}

	var docsCount int

	defer in.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(in)
	inStr := buf.String()
	lines := strings.FieldsFunc(inStr, func(r rune) bool {
		return r == '\n'
	})

	if c.Verbose {
		fmt.Printf("got %d lines\n", len(lines))
	}

	if len(lines) == 0 {
		return &model.ParseOutput{}
	}

	tokenIndex := make(map[string]*model.Tag)
	tokens := make([]string, 0)
	for _, l := range lines {
		sentences := util.SplitToSentences([]byte(l))
		for _, s := range sentences {
			docsCount++
			tokens = append(tokens, util.Sanitize(bytes.Fields(s), c.NoStopWords)...)
			visited := map[string]bool{}
			for _, token := range tokens {
				visited[token] = true
				item, ok := tokenIndex[token]
				if !ok {
					item = &model.Tag{Value: token}
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

	return &model.ParseOutput{Tags: tokenIndex, DocHash: hashTokens(tokens)}
}

func hashTokens(ts []string) []byte {
	h := sha512.New()
	for _, t := range ts {
		_, _ = h.Write([]byte(t))
	}
	return h.Sum(nil)
}
