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

func EmptyResult() *Result {
	return &Result{Meta: &Meta{}}
}

func ErrResult(err error) *Result {
	return &Result{Meta: &Meta{}, Err: err}
}

// Result represents result of Tagify.
type Result struct {
	Meta       *Meta
	RawTags    map[string]*Tag
	Tags       []*Tag
	Extensions map[string]map[string]*extension.Result
	Err        error
}

// FlatTags transforms internal token register into a slice.
func (res *Result) FlatTags() []*Tag {
	return flatten(res.RawTags)
}

// FindExtResults finds requested extension result(s), in case if version is empty.
func (res *Result) FindExtResults(name, version string) []*extension.Result {
	vs, ok := res.Extensions[name]
	if !ok {
		return nil
	}
	list := []*extension.Result{}
	if version != "" {
		if v, ok := vs[version]; ok {
			list = append(list, v)
		}
		return list
	}
	for _, v := range vs {
		list = append(list, v)
	}
	return list
}

// Meta extra information.
type Meta struct {
	ContentType config.ContentType
	DocTitle    string
	DocHash     string
	Lang        string
}

// Len returns count of tags in the result.
func (r *Result) Len() int {
	return len(r.RawTags)
}

// ForEach iterates through the slice of Tags
// and calls provided "fn" on every iteration.
func (r *Result) ForEach(fn func(i int, tag *Tag)) {
	for k, v := range r.Tags {
		fn(k, v)
	}
}

// TagsStrings transforms slice of tags into a slice of strings.
func (r *Result) TagsStrings() []string {
	return ToStrings(r.Tags)
}

// ParseFunc represents an arbitrary handler,
// which goes through given reader and produces tags.
type ParseFunc func(c *config.Config, reader io.ReadCloser) *Result

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
