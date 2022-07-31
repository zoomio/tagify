package config

import (
	"time"

	"github.com/zoomio/stopwords"

	"github.com/zoomio/tagify/extension"
)

var (
	allStopWords = map[string]stopwords.Option{
		"en": stopwords.Words(stopwords.StopWords),
		"ru": stopwords.Words(stopwords.StopWordsRu),
		"zh": stopwords.Words(stopwords.StopWordsZh),
		"ja": stopwords.Words(stopwords.StopWordsJa),
		"ko": stopwords.Words(stopwords.StopWordsKo),
		"hi": stopwords.Words(stopwords.StopWordsHi),
		"he": stopwords.Words(stopwords.StopWordsHe),
		"ar": stopwords.Words(stopwords.StopWordsAr),
		"de": stopwords.Words(stopwords.StopWordsDe),
		"es": stopwords.Words(stopwords.StopWordsEs),
		"fr": stopwords.Words(stopwords.StopWordsFr),
	}
)

// New ...
func New(options ...Option) *Config {
	c := &Config{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	return c
}

// Config ...
type Config struct {
	Source string

	// headless
	Query      string
	WaitFor    string
	WaitUntil  time.Duration
	Screenshot bool

	Content string
	ContentType
	Limit         int
	Verbose       bool
	NoStopWords   bool
	SkipLang      bool
	Lang          string
	StopWords     *stopwords.Register
	ContentOnly   bool
	FullSite      bool
	AllTagWeights bool
	TagWeights
	ExtraTagWeights TagWeights
	ExcludeTags     TagWeights
	AdjustScores    bool
	Extensions      []extension.Extension
}

// SetStopWords ...
func (c *Config) SetStopWords(lang string) {
	c.Lang = lang
	if found, ok := allStopWords[lang]; ok {
		c.StopWords = stopwords.Setup(found)
	} else {
		c.StopWords = stopwords.Setup(stopwords.Words(stopwords.StopWords))
	}
}
