package tagify

import (
	"context"
	"strings"

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
func (ct ContentType) String() string {
	if ct < Text || ct > HTML {
		return "Unknown"
	}
	return contentTypes[ct]
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

	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
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

// Read reads into given bytes (does not close reader).
// Makes `in` to be compatible with `io.Reader`.
func (in *in) Read(p []byte) (n int, err error) {
	return in.reader.Read(p)
}

// Close closes internal reader.
// Makes `in` to be compatible with `io.Closer`.
func (in *in) Close() error {
	return in.reader.Close()
}

// ReadLines provides slice of lines from input,
// after method has returned input is closed.
func (in *in) ReadLines() ([]string, error) {
	lines, err := in.reader.ReadLines()
	if err != nil {
		return nil, err
	}
	return lines, nil
}
