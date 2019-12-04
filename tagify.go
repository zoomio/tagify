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

// GetTags produces slice of tags ordered by frequency and limited by limit.
//
// Deprecated: use tagify#Run instead.
func GetTags(ctx context.Context, source string, contentType ContentType, limit int, verbose, noStopWords bool) ([]*processor.Tag, error) {
	return Run(ctx, Source(source), TargetType(contentType),
		Limit(limit), Verbose(verbose), NoStopWords(noStopWords))
}

// GetTagsWithQuery produces slice of tags from "source" narrowed down to a CSS "query" ordered by frequency and limited by limit.
//
// Deprecated: use tagify#Run instead.
func GetTagsWithQuery(ctx context.Context, source, query string, contentType ContentType, limit int,
	verbose, noStopWords bool) ([]*processor.Tag, error) {
	return Run(ctx, Source(source), Query(query), TargetType(contentType),
		Limit(limit), Verbose(verbose), NoStopWords(noStopWords))
}

// Run produces slice of tags ordered by frequency.
func Run(ctx context.Context, options ...Option) ([]*processor.Tag, error) {

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

	return processInput(&in, *c)
}

// GetTagsFromString produces slice of tags ordered by frequency and limited by limit.
func GetTagsFromString(input string, contentType ContentType, limit int, verbose, noStopWords bool) ([]*processor.Tag, error) {
	in := newInFromString(input, contentType)

	c := &config{
		limit:       limit,
		verbose:     verbose,
		noStopWords: noStopWords,
	}

	return processInput(&in, *c)
}

// ToStrings transforms a list of tags into a list of strings.
func ToStrings(items []*processor.Tag) []string {
	return processor.ToStrings(items)
}

func processInput(in *in, c config) ([]*processor.Tag, error) {
	var tags []*processor.Tag

	if c.verbose {
		fmt.Println("reading lines...")
	}

	switch in.ContentType {
	case HTML:
		if c.verbose {
			fmt.Println("parsing HTML...")
		}
		tags = processor.ParseHTML(in.getReader(), c.verbose, c.noStopWords)
	default:
		lines, err := in.readAllLines()
		if err != nil {
			return tags, err
		}

		if c.verbose {
			fmt.Printf("got %d lines\n", len(lines))
		}

		if len(lines) == 0 {
			return tags, nil
		}

		if c.verbose {
			fmt.Println("parsing plain text...")
		}
		tags = processor.ParseText(lines, c.noStopWords)
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

	return tags, nil
}
