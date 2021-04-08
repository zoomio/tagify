package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseOption allows to customise `Tagger` configuration.
type ParseOption func(*ParseConfig)

var (
	// Verbose enables high verbosity.
	Verbose = func(verbose bool) ParseOption {
		return func(c *ParseConfig) {
			c.Verbose = verbose
		}
	}

	// NoStopWords enables stop-words exclusion from the output.
	NoStopWords = func(noStopWords bool) ParseOption {
		return func(c *ParseConfig) {
			c.NoStopWords = noStopWords
		}
	}

	// ContentOnly ignores all none content related parts of the HTML page (HTML only).
	ContentOnly = func(v bool) ParseOption {
		return func(c *ParseConfig) {
			c.ContentOnly = v
		}
	}

	// FullSite tells parser to process full site (HTML only).
	FullSite = func(v bool) ParseOption {
		return func(c *ParseConfig) {
			c.FullSite = v
		}
	}

	// Source of the parser.
	Source = func(v string) ParseOption {
		return func(c *ParseConfig) {
			c.Source = v
		}
	}

	TagWeightsString = func(v string) ParseOption {
		return func(c *ParseConfig) {
			c.TagWeights = ParseTagWeights(strings.NewReader(v), String)
		}
	}

	TagWeightsJSON = func(v string) ParseOption {
		return func(c *ParseConfig) {
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
)
