package tagify

import (
	"github.com/zoomio/tagify/processor/model"
)

// Result represents result of Tagify.
type Result struct {
	Meta *Meta
	Tags []*model.Tag
}

// Meta extra information.
type Meta struct {
	ContentType ContentType
	DocTitle    string
	DocHash     string
}

// Len returns count of tags in the result.
func (r *Result) Len() int {
	return len(r.Tags)
}

// ForEach iterates through the slice of Tags
// and calls provided "fn" on every iteration.
func (r *Result) ForEach(fn func(i int, tag *model.Tag)) {
	for k, v := range r.Tags {
		fn(k, v)
	}
}

// TagsStrings transforms slice of tags into a slice of strings.
func (r *Result) TagsStrings() []string {
	return ToStrings(r.Tags)
}
