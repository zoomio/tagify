package tagify

// Option allows to customise configuration.
type Option func(*config)

var (
	// Source sets target source.
	Source = func(source string) Option {
		return func(c *config) {
			c.source = source
		}
	}

	// Query sets CSS query for the target.
	Query = func(query string) Option {
		return func(c *config) {
			c.query = query
		}
	}

	// Content sets content of the target.
	Content = func(content string) Option {
		return func(c *config) {
			c.content = content
		}
	}

	// TargetType sets content type of the target.
	TargetType = func(contentType ContentType) Option {
		return func(c *config) {
			c.contentType = contentType
		}
	}

	// Limit sets cthe limit of tags for the target.
	Limit = func(limit int) Option {
		return func(c *config) {
			c.limit = limit
		}
	}

	// Verbose enables high verbosity.
	Verbose = func(verbose bool) Option {
		return func(c *config) {
			c.verbose = verbose
		}
	}

	// NoStopWords enables stop-words exclusion from the output.
	NoStopWords = func(noStopWords bool) Option {
		return func(c *config) {
			c.noStopWords = noStopWords
		}
	}
)
