package config

// Content types
const (
	Unknown ContentType = iota
	Text
	HTML
	Markdown
)

var (
	ContentTypes = [...]string{
		"Unknown",
		"Text",
		"HTML",
		"Markdown",
	}
)

// ContentType ...
type ContentType byte

// ContentTypeOf returns ContentType based on string value.
func ContentTypeOf(contentType string) ContentType {
	for i, key := range ContentTypes {
		if key == contentType {
			return ContentType(i)
		}
	}
	return Unknown
}

// String ...
func (ct ContentType) String() string {
	if ct < Text || ct > Markdown {
		return "Unknown"
	}
	return ContentTypes[ct]
}
