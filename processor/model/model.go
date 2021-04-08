package model

import (
	"fmt"
	"io"
)

// InputReader ...
type InputReader interface {
	ReadLines() ([]string, error)
	io.ReadCloser
}

// Tag holds some arbitrary string value (e.g. a word) along with some extra data about it.
type Tag struct {
	// Value of the tag, i.e. a word
	Value string
	// Score used to represent importance of the tag
	Score float64
	// Count is the number of times tag appeared in a text
	Count int
	// Docs is the number of documents in a text in which the tag appeared
	Docs int
	// DocsCount is the number of documents in a text
	DocsCount int
}

func (t *Tag) String() string {
	return fmt.Sprintf("(%s - [score: %.2f, count: %d, docs: %d, docs_count: %d])",
		t.Value, t.Score, t.Count, t.Docs, t.DocsCount)
}

// ParseOutput is a result of the `ParseFunc`.
type ParseOutput struct {
	Tags     map[string]*Tag
	DocTitle string
	DocHash  []byte
	Err      error
}

// FlatTags transforms internal token register into a slice.
func (po *ParseOutput) FlatTags() []*Tag {
	return flatten(po.Tags)
}

type ParseConfig struct {
	Verbose     bool
	NoStopWords bool
	ContentOnly bool
	FullSite    bool
	Source      string
}

// ParseFunc represents an arbitrary handler,
// which goes through given reader and produces tags.
type ParseFunc func(reader io.ReadCloser, options ...ParseOption) *ParseOutput

func flatten(dict map[string]*Tag) []*Tag {
	flat := make([]*Tag, len(dict))
	var i int
	for _, val := range dict {
		flat[i] = val
		i++
	}
	return flat
}

// ToStrings transforms list of given tags into a list of strings.
func ToStrings(items []*Tag) []string {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = item.Value
	}
	return strs
}
