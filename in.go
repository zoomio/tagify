package tagify

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/zoomio/inout"
)

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
func newIn(ctx context.Context, cfg *Config) (in, error) {
	in := in{source: cfg.Source}

	if strings.HasPrefix(cfg.Source, "http://") || strings.HasPrefix(cfg.Source, "https://") || cfg.Query != "" {
		in.ContentType = HTML
	} else if strings.ToLower(filepath.Ext(cfg.Source)) == ".md" {
		in.ContentType = Markdown
	}

	r, err := inout.NewInOut(ctx,
		inout.Source(cfg.Source),
		inout.Query(cfg.Query),
		inout.WaitFor(cfg.WaitFor),
		inout.WaitUntil(cfg.WaitUntil),
		inout.Screenshot(cfg.Screenshot),
		inout.Timeout(cfg.Timeout),
		inout.Verbose(cfg.Verbose))
	if err != nil {
		return in, err
	}

	in.reader = &r

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
