package tagify

// Option allows to customise configuration.
type Option func(*config)

var (
	// Source sets target source.
	Source = func(v string) Option {
		return func(c *config) {
			c.source = v
		}
	}

	// Query sets CSS query for the target.
	Query = func(v string) Option {
		return func(c *config) {
			c.query = v
		}
	}

	// Content sets content of the target.
	Content = func(v string) Option {
		return func(c *config) {
			c.content = v
		}
	}

	// TargetType sets content type of the target.
	TargetType = func(v ContentType) Option {
		return func(c *config) {
			c.contentType = v
		}
	}

	// Limit sets the limit of tags for the target.
	Limit = func(v int) Option {
		return func(c *config) {
			c.limit = v
		}
	}

	// Verbose enables high verbosity.
	Verbose = func(v bool) Option {
		return func(c *config) {
			c.verbose = v
		}
	}

	// NoStopWords enables stop-words exclusion from the output.
	NoStopWords = func(v bool) Option {
		return func(c *config) {
			c.noStopWords = v
		}
	}

	// ContentOnly ignores all none content related parts of the HTML page.
	ContentOnly = func(v bool) Option {
		return func(c *config) {
			c.contentOnly = v
		}
	}

	// FullSite tells parser to process full site (HTML only).
	FullSite = func(v bool) Option {
		return func(c *config) {
			c.fullSite = v
		}
	}

	// TagWeights string with the custom tag weights for the HTML & Markdown tagging.
	TagWeights = func(v string) Option {
		return func(c *config) {
			c.tagWeights = v
		}
	}

	// TagWeightsJSON JSON file with the custom tag weights for the HTML & Markdown tagging.
	TagWeightsJSON = func(v string) Option {
		return func(c *config) {
			c.tagWeightsJSON = v
		}
	}
)
