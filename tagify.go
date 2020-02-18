package tagify

import (
	"context"
	"fmt"

	"github.com/zoomio/tagify/processor"
)

type config struct {
	source      string
	query       string
	content     string
	contentType ContentType
	limit       int
	verbose     bool
	noStopWords bool
}

// Run produces slice of tags ordered by frequency.
func Run(ctx context.Context, options ...Option) (*Result, error) {

	c := &config{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	var in in
	var err error

	if c.content != "" {
		in = newInFromString(c.content, c.contentType)
	} else {
		in, err = newIn(ctx, c.source, c.query, c.verbose)
		if c.contentType > Unknown {
			in.ContentType = c.contentType
		} else if c.query != "" {
			in.ContentType = HTML
		}
	}

	if err != nil {
		return nil, err
	}

	tags, title, version := processInput(&in, *c)

	return &Result{
		Meta: &Meta{
			ContentType: in.ContentType,
			DocTitle:    title,
			DocVersion:  fmt.Sprintf("%x", version),
		},
		Tags: tags,
	}, nil
}

// ToStrings transforms a list of tags into a list of strings.
func ToStrings(items []*processor.Tag) []string {
	return processor.ToStrings(items)
}

func processInput(in *in, c config) (tags []*processor.Tag, pageTitle string, version []byte) {
	switch in.ContentType {
	case HTML:
		tags, pageTitle, version = processor.ParseHTML(in, c.verbose, c.noStopWords)
	default:
		tags, version = processor.ParseText(in, c.verbose, c.noStopWords)
	}

	if len(tags) > 0 {
		if c.verbose {
			fmt.Println("tagifying...")
		}
		tags = processor.Run(tags, c.limit)
		if c.verbose {
			fmt.Printf("\n%v\n", tags)
		}
	}

	return
}
