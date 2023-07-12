package tagify

import (
	"github.com/zoomio/tagify/config"
)

// backwards compatibility
type Config = config.Config
type Option = config.Option
type ContentType = config.ContentType

var (
	Source   = config.Source
	Language = config.Language
	Content  = config.Content

	Timeout = config.Timeout

	// headless
	Query      = config.Query
	WaitFor    = config.WaitFor
	WaitUntil  = config.WaitUntil
	Screenshot = config.Screenshot

	// misc
	TargetType  = config.TargetType
	Limit       = config.Limit
	Verbose     = config.Verbose
	NoStopWords = config.NoStopWords
	StopWords   = config.StopWords
	ContentOnly = config.ContentOnly
	FullSite    = config.FullSite

	// weighing
	TagWeightsString      = config.TagWeightsString
	TagWeightsJSON        = config.TagWeightsJSON
	ExtraTagWeightsString = config.ExtraTagWeightsString
	ExtraTagWeightsJSON   = config.ExtraTagWeightsJSON
	ExcludeTagsString     = config.ExcludeTagsString
	AllTagWeights         = config.AllTagWeights
	AdjustScores          = config.AdjustScores

	// content types
	Unknown       = config.Unknown
	Text          = config.Text
	HTML          = config.HTML
	Markdown      = config.Markdown
	ContentTypeOf = config.ContentTypeOf

	Extensions = config.Extensions
)
