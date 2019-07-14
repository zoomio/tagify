package tagify

import (
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
func GetTags(source string, contentType ContentType, limit int, verbose, noStopWords bool) ([]*processor.Tag, error) {
	return Run(Source(source), TargetType(contentType),
		Limit(limit), Verbose(verbose), NoStopWords(noStopWords))
}

// GetTagsWithQuery produces slice of tags from "source" narrowed down to a CSS "query" ordered by frequency and limited by limit.
//
// Deprecated: use tagify#Run instead.
func GetTagsWithQuery(source, query string, contentType ContentType, limit int,
	verbose, noStopWords bool) ([]*processor.Tag, error) {
	return Run(Source(source), Query(query), TargetType(contentType),
		Limit(limit), Verbose(verbose), NoStopWords(noStopWords))
}

// Run produces slice of tags ordered by frequency.
func Run(options ...Option) ([]*processor.Tag, error) {

	c := &config{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	in, err := NewIn(c.source, c.query)
	if err != nil {
		return []*processor.Tag{}, err
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
	in := NewInFromString(input, contentType)

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

func processInput(in *In, c config) ([]*processor.Tag, error) {
	var tags []*processor.Tag

	lines, err := in.ReadAllLines()
	if err != nil {
		return tags, err
	}

	switch in.ContentType {
	case HTML:
		tags = processor.ParseHTML(lines, c.verbose, c.noStopWords)
	default:
		tags = processor.ParseText(lines, c.noStopWords)
	}

	tags = processor.Run(tags, c.limit)
	if c.verbose {
		fmt.Printf("%v\n", tags)
		fmt.Printf("\nsize: %d\n\n", len(tags))
	}

	return tags, nil
}
