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

	// ContentOnly ignores all none content related parts of the HTML page (HTML only).
	ContentOnly = func(v bool) ParseOption {
		return func(c *parseConfig) {
			c.contentOnly = v
		}
	}

	// FullSite tells parser to process full site (HTML only).
	FullSite = func(v bool) ParseOption {
		return func(c *parseConfig) {
			c.fullSite = v
		}
	}

	// Source of the parser.
	Source = func(v string) ParseOption {
		return func(c *parseConfig) {
			c.source = v
		}
	}
)
