package model

import (
	"fmt"
	"io"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
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
	Tags       map[string]*Tag
	DocTitle   string
	DocHash    []byte
	Lang       string
	Err        error
	Extensions map[string]map[string]*extension.Result
}

// FlatTags transforms internal token register into a slice.
func (po *ParseOutput) FlatTags() []*Tag {
	return flatten(po.Tags)
}

// FindExtResults finds requested extension result(s), in case if version is empty.
func (c *ParseOutput) FindExtResults(name, version string) []*extension.Result {
	vs, ok := c.Extensions[name]
	if !ok {
		return nil
	}
	res := []*extension.Result{}
	if version != "" {
		if v, ok := vs[version]; ok {
			res = append(res, v)
		}
		return res
	}
	for _, v := range vs {
		res = append(res, v)
	}
	return res
}

// ParseFunc represents an arbitrary handler,
// which goes through given reader and produces tags.
type ParseFunc func(c *config.Config, reader io.ReadCloser) *ParseOutput

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
