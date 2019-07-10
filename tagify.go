package tagify

import (
	"fmt"

	"github.com/zoomio/tagify/processor"
)

func processInput(in *In, limit int, verbose, noStopWords bool) ([]*processor.Tag, error) {
	var tags []*processor.Tag

	lines, err := in.ReadAllLines()
	if err != nil {
		return tags, err
	}

	switch in.ContentType {
	case HTML:
		tags = processor.ParseHTML(lines, verbose, noStopWords)
	default:
		tags = processor.ParseText(lines, noStopWords)
	}

	tags = processor.Run(tags, limit)
	if verbose {
		fmt.Printf("%v\n", tags)
		fmt.Printf("\nsize: %d\n\n", len(tags))
	}

	return tags, nil
}

// GetTags produces slice of tags ordered by frequency and limited by limit.
func GetTags(source string, contentType ContentType, limit int, verbose, noStopWords bool) ([]*processor.Tag, error) {
	return GetTagsWithQuery(source, "", contentType, limit, verbose, noStopWords)
}

// GetTagsWithQuery produces slice of tags from "source" narrowed down to a CSS "query" ordered by frequency and limited by limit.
func GetTagsWithQuery(source, query string, contentType ContentType, limit int,
	verbose, noStopWords bool) ([]*processor.Tag, error) {
	in, err := NewIn(source, query)
	if err != nil {
		return []*processor.Tag{}, err
	}
	if query != "" {
		in.ContentType = Text
	}
	if contentType > Unknown {
		in.ContentType = contentType
	}

	return processInput(&in, limit, verbose, noStopWords)
}

// GetTagsFromString produces slice of tags ordered by frequency and limited by limit.
func GetTagsFromString(input string, contentType ContentType, limit int, verbose, noStopWords bool) ([]*processor.Tag, error) {
	in := NewInFromString(input, contentType)
	return processInput(&in, limit, verbose, noStopWords)
}

// ToStrings transforms a list of tags into a list of strings.
func ToStrings(items []*processor.Tag) []string {
	return processor.ToStrings(items)
}
