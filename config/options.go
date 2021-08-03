package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	// Query sets CSS query for the target.
	Query = func(v string) Option {
		return func(c *Config) {
			c.Query = v
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

	AdjustScores = func(v bool) Option {
		return func(c *Config) {
			c.AdjustScores = v
		}
	}
)
