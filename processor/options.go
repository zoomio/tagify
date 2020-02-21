package processor

// ParseOption allows to customise `Tagger` configuration.
type ParseOption func(*parseConfig)

var (
	// Verbose enables high verbosity.
	Verbose = func(verbose bool) ParseOption {
		return func(c *parseConfig) {
			c.verbose = verbose
		}
	}

	// NoStopWords enables stop-words exclusion from the output.
	NoStopWords = func(noStopWords bool) ParseOption {
		return func(c *parseConfig) {
			c.noStopWords = noStopWords
		}
	}
)
