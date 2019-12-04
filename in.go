package tagify

import (
	"context"
	"io"
	"os"

	"github.com/zoomio/inout"
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
type ContentType byte

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

// in - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type in struct {
	reader *inout.Reader
	ContentType
}

// newIn initializes an input stream from STDIN, file or web page.
//
// name - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func newIn(ctx context.Context, name, query string, verbose bool) (in, error) {
	in := in{}
	r, err := inout.NewInOut(ctx, inout.Source(name), inout.Query(query), inout.Verbose(verbose))
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

// newInFromString ...
func newInFromString(input string, contentType ContentType) in {
	r := inout.NewFromString(input)
	return in{
		ContentType: contentType,
		reader:      r,
	}
}

func (in *in) getReader() io.Reader {
	return in.reader
}

// readAllStrings provides slice of strings from input split by white space.
func (in *in) readAllStrings() ([]string, error) {
	strs, err := in.reader.ReadWords()
	if err != nil {
		return nil, err
	}
	return strs, nil
}

// readAllLines provides slice of lines from input.
func (in *in) readAllLines() ([]string, error) {
	lines, err := in.reader.ReadLines()
	if err != nil {
		return nil, err
	}
	return lines, nil
}
