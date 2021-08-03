package config

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

	// TagWeights string with the custom tag weights for the HTML & Markdown tagging.
	TagWeightsStr = func(v string) Option {
		return func(c *Config) {
			c.TagWeightsStr = v
		}
	}

	// TagWeightsJSON JSON file with the custom tag weights for the HTML & Markdown tagging.
	TagWeightsJSON = func(v string) Option {
		return func(c *Config) {
			c.TagWeightsJSON = v
		}
	}

	AdjustScores = func(v bool) Option {
		return func(c *Config) {
			c.AdjustScores = v
		}
	}
)
