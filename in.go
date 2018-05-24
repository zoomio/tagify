package tagify

import (
	"os"

	"github.com/zoomio/tagify/reader"
)

// Content types
const (
	Unknown ContentType = iota
	Text
	HTML
)

var (
	contentTypes = [...]string{
		"Unknown",
		"Text",
		"HTML",
	}
)

// ContentType ...
type ContentType int

// ContentTypeOf returns ContentType based on string value.
func ContentTypeOf(contentType string) ContentType {
	for i, key := range contentTypes {
		if key == contentType {
			return ContentType(i)
		}
	}
	return Unknown
}

// String ...
func (contentType ContentType) String() string {
	if contentType < Text || contentType > HTML {
		return "Unknown"
	}
	return contentTypes[contentType]
}

// In - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type In struct {
	reader *reader.Reader
	ContentType
}

// NewIn initializes an input stream from STDIN, file or web page.
//
// name - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func NewIn(name string) (In, error) {
	in := In{}
	r, err := reader.NewIn(name)
	if err != nil {
		return in, err
	}

	in.reader = &r

	_, statErr := os.Stat(name)
	if name != "" && statErr != nil {
		in.ContentType = HTML
	}

	return in, err
}

// NewInFromString ...
func NewInFromString(input string, contentType ContentType) In {
	r := reader.NewInFromString(input)
	return In{
		ContentType: contentType,
		reader:      &r,
	}
}

// ReadAllStrings provides slice of strings from input split by white space.
func (in *In) ReadAllStrings() ([]string, error) {
	strs, err := in.reader.ReadAllStrings()
	if err != nil {
		return []string{}, err
	}
	return strs, nil
}

// ReadAllLines provides slice of lines from input.
func (in *In) ReadAllLines() ([]string, error) {
	lines, err := in.reader.LinesFromReader()
	if err != nil {
		return []string{}, err
	}
	return lines, nil
}
