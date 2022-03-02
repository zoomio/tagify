package tagify

import (
	"github.com/zoomio/tagify/config"
)

// backwards compatibility
type Config = config.Config
type Option = config.Option
type ContentType = config.ContentType

var (
	Source                = config.Source
	Query                 = config.Query
	Content               = config.Content
	TargetType            = config.TargetType
	Limit                 = config.Limit
	Verbose               = config.Verbose
	NoStopWords           = config.NoStopWords
	StopWords             = config.StopWords
	ContentOnly           = config.ContentOnly
	FullSite              = config.FullSite
	TagWeightsString      = config.TagWeightsString
	TagWeightsJSON        = config.TagWeightsJSON
	ExtraTagWeightsString = config.ExtraTagWeightsString
	ExtraTagWeightsJSON   = config.ExtraTagWeightsJSON
	AdjustScores          = config.AdjustScores
	Extensions            = config.Extensions

	Unknown  = config.Unknown
	Text     = config.Text
	HTML     = config.HTML
	Markdown = config.Markdown

	ContentTypeOf = config.ContentTypeOf
)
