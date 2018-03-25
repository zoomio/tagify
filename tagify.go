package tagify

import (
	"github.com/zoomio/tagify/rank"
	"github.com/zoomio/tagify/inout"
)

func toStrings(items []*rank.Item) []string {
	strs := make([]string, len(items))
	var i int
	for _, item := range items {
		strs[i] = item.Value
		i++
	}
	return strs
}

func processInput(in *inout.In, limit int, verbose bool) ([]string, error) {
	var items []*rank.Item

	switch in.ContentType {
	case inout.HTML:
		items = rank.ParseHTML(in.GetLines(), verbose)		
	default:
		items = rank.ParseText(in.ReadAllStrings())
	}

	sortByScoreDescending(items)
	if limit > 0 {
		return toStrings(items[:limit]), nil
	}
	return toStrings(items), nil
}

// GetTags produces slice of tags ordered by frequency and limited by limit.
func GetTags(source string, contentType, limit int, verbose bool) ([]string, error) {
	in, err := inout.NewIn(source, contentType)
	if err != nil {
		return []string{}, err
	}
	return processInput(&in, limit, verbose)
}

// GetTagsFromString produces slice of tags ordered by frequency and limited by limit.
func GetTagsFromString(input string, contentType, limit int, verbose bool) ([]string, error) {
	in, err := inout.NewInFromString(input, contentType)
	if err != nil {
		return []string{}, err
	}
	return processInput(&in, limit, verbose)
}