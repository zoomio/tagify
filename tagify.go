package tagify

import (
	"context"
	"fmt"

	"github.com/zoomio/tagify/processor"
)

type config struct {
	source      string
	query       string
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

	in, err := newIn(ctx, c.source, c.query, c.verbose)
	if err != nil {
		return nil, err
	}
	if c.query != "" {
		in.ContentType = Text
	}
	if c.contentType > Unknown {
		in.ContentType = c.contentType
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

// GetTagsFromString produces slice of tags ordered by frequency and limited by limit.
func GetTagsFromString(input string, contentType ContentType, limit int, verbose, noStopWords bool) ([]*processor.Tag, []byte) {
	in := newInFromString(input, contentType)

	c := &config{
		limit:       limit,
		verbose:     verbose,
		noStopWords: noStopWords,
	}

	tags, _, version := processInput(&in, *c)

	return tags, version
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
