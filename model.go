package tagify

import (
	"github.com/zoomio/tagify/processor"
)

// Result represents result of Tagify.
type Result struct {
	Meta *Meta
	Tags []*processor.Tag
}

// Meta extra information.
type Meta struct {
	ContentType ContentType
	DocTitle    string
}

// Len returns count of tags in the result.
func (r *Result) Len() int {
	return len(r.Tags)
}

// TagsStrings transforms slice of tags into a slice of strings.
func (r *Result) TagsStrings() []string {
	return ToStrings(r.Tags)
}
