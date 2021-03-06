package tagify

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/zoomio/inout"
)

// Content types
const (
	Unknown ContentType = iota
	Text
	HTML
	Markdown
)

var (
	contentTypes = [...]string{
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
	for i, key := range contentTypes {
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
	return contentTypes[ct]
}

// in - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type in struct {
	source string
	reader *inout.Reader
	ContentType
}

// newIn initializes an input stream from STDIN, file or web page.
//
// source - the filename or web page source, reads from STDIN if source is empty.
// Panics on errors.
func newIn(ctx context.Context, source, query string, verbose bool) (in, error) {
	in := in{source: source}
	r, err := inout.NewInOut(ctx,
		inout.Source(source),
		inout.Query(query),
		inout.Timeout(0),
		inout.Verbose(verbose))
	if err != nil {
		return in, err
	}

	in.reader = &r

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") || query != "" {
		in.ContentType = HTML
	} else if strings.ToLower(filepath.Ext(source)) == ".md" {
		in.ContentType = Markdown
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
