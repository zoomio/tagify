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
	PageTitle   string
	ContentType ContentType
}

// Len ...
func (r *Result) Len() int {
	return len(r.Tags)
}

// TagsStrings ...
func (r *Result) TagsStrings() []string {
	return ToStrings(r.Tags)
}
