package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zoomio/stopwords"
	"github.com/zoomio/tagify/extension"
)

// Option allows to customise configuration.
type Option func(*Config)

var (
	// Source sets target source.
	Source = func(v string) Option {
		return func(c *Config) {
			c.Source = v
		}
	}

	// Language ...
	Language = func(v string) Option {
		return func(c *Config) {
			c.Lang = v
		}
	}

	// Query sets CSS query for the target.
	Query = func(v string) Option {
		return func(c *Config) {
			c.Query = v
		}
	}

	// Timeout sets the overall deadline for the operation.
	Timeout = func(d time.Duration) Option {
		return func(c *Config) {
			c.Timeout = d
		}
	}

	// WaitFor sets CSS query for the target of In-Out.
	WaitFor = func(query string) Option {
		return func(c *Config) {
			c.WaitFor = query
		}
	}

	// WaitUntil sets page load duration to wait for.
	WaitUntil = func(d time.Duration) Option {
		return func(c *Config) {
			c.WaitUntil = d
		}
	}

	// Screenshot captures screenshot, Reader will ImgBytes of the image populated.
	Screenshot = func(v bool) Option {
		return func(c *Config) {
			c.Screenshot = v
		}
	}

	// Content sets content of the target.
	Content = func(v string) Option {
		return func(c *Config) {
			c.Content = v
		}
	}

	// TargetType sets content type of the target.
	TargetType = func(v ContentType) Option {
		return func(c *Config) {
			c.ContentType = v
		}
	}

	// Limit sets the limit of tags for the target.
	Limit = func(v int) Option {
		return func(c *Config) {
			c.Limit = v
		}
	}

	// Verbose enables high verbosity.
	Verbose = func(v bool) Option {
		return func(c *Config) {
			c.Verbose = v
		}
	}

	// NoStopWords enables stop-words exclusion from the output.
	NoStopWords = func(v bool) Option {
		return func(c *Config) {
			c.NoStopWords = v
		}
	}

	// StopWords allows to provide a custom set of stop-words.
	StopWords = func(v []string) Option {
		return func(c *Config) {
			c.StopWords = stopwords.Setup(stopwords.WordsSlice(v))
		}
	}

	// ContentOnly ignores all none content related parts of the HTML page.
	ContentOnly = func(v bool) Option {
		return func(c *Config) {
			c.ContentOnly = v
		}
	}

	// FullSite tells parser to process full site (HTML only).
	FullSite = func(v bool) Option {
		return func(c *Config) {
			c.FullSite = v
		}
	}

	// TagWeightsString ...
	TagWeightsString = func(v string) Option {
		return func(c *Config) {
			c.TagWeights = ParseTagWeights(strings.NewReader(v), String)
		}
	}

	// TagWeightsJSON ...
	TagWeightsJSON = func(v string) Option {
		return func(c *Config) {
			f, err := os.Open(v)
			if err != nil {
				println(fmt.Errorf("error: can't open JSON file [%s]: %w", v, err))
				return
			}
			r := bufio.NewReader(f)
			c.TagWeights = ParseTagWeights(r, JSON)
			f.Close()
		}
	}

	// ExtraTagWeightsString ...
	ExtraTagWeightsString = func(v string) Option {
		return func(c *Config) {
			c.ExtraTagWeights = ParseTagWeights(strings.NewReader(v), String)
		}
	}

	// TagWeightsJSON ...
	ExtraTagWeightsJSON = func(v string) Option {
		return func(c *Config) {
			f, err := os.Open(v)
			if err != nil {
				println(fmt.Errorf("error: can't open JSON file [%s]: %w", v, err))
				return
			}
			r := bufio.NewReader(f)
			c.ExtraTagWeights = ParseTagWeights(r, JSON)
			f.Close()
		}
	}

	// ExcludeTagsString ...
	ExcludeTagsString = func(v string) Option {
		return func(c *Config) {
			c.ExcludeTags = ParseTagWeights(strings.NewReader(v), String)
		}
	}

	// AllTagWeights ...
	AllTagWeights = func(v bool) Option {
		return func(c *Config) {
			c.AllTagWeights = v
		}
	}

	AdjustScores = func(v bool) Option {
		return func(c *Config) {
			c.AdjustScores = v
		}
	}

	Extensions = func(v []extension.Extension) Option {
		return func(c *Config) {
			c.Extensions = make([]extension.Extension, len(v))
			copy(c.Extensions, v)
		}
	}
)
