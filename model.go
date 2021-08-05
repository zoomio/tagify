package tagify

import (
	"github.com/zoomio/tagify/config"
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

// backwards compatibility
type Config = config.Config
type Option = config.Option
type ContentType = config.ContentType

var (
	Source           = config.Source
	Query            = config.Query
	Content          = config.Content
	TargetType       = config.TargetType
	Limit            = config.Limit
	Verbose          = config.Verbose
	NoStopWords      = config.NoStopWords
	ContentOnly      = config.ContentOnly
	FullSite         = config.FullSite
	TagWeightsString = config.TagWeightsString
	TagWeightsJSON   = config.TagWeightsJSON
	AdjustScores     = config.AdjustScores
	Extensions       = config.Extensions

	Unknown  = config.Unknown
	Text     = config.Text
	HTML     = config.HTML
	Markdown = config.Markdown

	ContentTypeOf = config.ContentTypeOf
)
