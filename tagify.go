package tagify

import (
	"fmt"

	"github.com/zoomio/tagify/processor"
)

func processInput(in *In, limit int, verbose, doFiltering bool) ([]*processor.Tag, error) {
	var tags []*processor.Tag

	lines, err := in.ReadAllLines()
	if err != nil {
		return tags, err
	}

	switch in.ContentType {
	case HTML:
		_, tags = processor.ParseHTML(lines, true, verbose, doFiltering)
	default:
		tags = processor.ParseText(lines, doFiltering)
	}

	tags = processor.Run(tags, limit)
	if verbose {
		fmt.Printf("%v\n", tags)
		fmt.Printf("\nsize: %d\n\n", len(tags))
	}

	return tags, nil
}

// GetTags produces slice of tags ordered by frequency and limited by limit.
func GetTags(source string, contentType ContentType, limit int, verbose, filterStopwords bool) ([]*processor.Tag, error) {
	in, err := NewIn(source)
	if err != nil {
		return []*processor.Tag{}, err
	}
	if contentType > Unknown {
		in.ContentType = contentType
	}

	return processInput(&in, limit, verbose, filterStopwords)
}

// GetTagsFromString produces slice of tags ordered by frequency and limited by limit.
func GetTagsFromString(input string, contentType ContentType, limit int, verbose, filterStopwords bool) ([]*processor.Tag, error) {
	in := NewInFromString(input, contentType)
	return processInput(&in, limit, verbose, filterStopwords)
}

// ToStrings ...
func ToStrings(items []*processor.Tag) []string {
	return processor.ToStrings(items)
}
