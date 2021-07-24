package config

import "github.com/zoomio/stopwords"

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

// Config ...
type Config struct {
	Source  string
	Query   string
	Content string
	ContentType
	Limit          int
	Verbose        bool
	NoStopWords    bool
	Lang           string
	StopWords      *stopwords.Register
	ContentOnly    bool
	FullSite       bool
	TagWeights     string
	TagWeightsJSON string
	AdjustScores   bool
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

// New ...
func New(options ...Option) *Config {
	c := &Config{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	return c
}
