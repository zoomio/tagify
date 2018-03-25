package tagify

import (
	"fmt"
	"github.com/gobuffalo/packr"

	"github.com/zoomio/tagify/rank"
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

func processInput(in *In, limit int, verbose bool) ([]string, error) {
	var items []*rank.Item

	switch in.ContentType {
	case HTML:
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

// Init initializes Tagify.
func Init() error {
	box := packr.NewBox("./_files")
	in, err := NewInFromString(box.String("stop-word-list.txt"), Text)
	if err != nil {
		return fmt.Errorf("error in initialization of Tagify: %v", err)
	}
	rank.RegisterStopWords(in.ReadAllStrings())
	return nil
}

// GetTags produces slice of tags ordered by frequency and limited by limit.
func GetTags(source, contentType string, limit int, verbose bool) ([]string, error) {
	in, err := NewIn(source)
	if err != nil {
		return []string{}, err
	}
	t := ContentTypeOf(contentType)
	if t > Unknown {
		in.ContentType = t
	}
	return processInput(&in, limit, verbose)
}

// GetTagsFromString produces slice of tags ordered by frequency and limited by limit.
func GetTagsFromString(input, contentType string, limit int, verbose bool) ([]string, error) {
	in, err := NewInFromString(input, ContentTypeOf(contentType))
	if err != nil {
		return []string{}, err
	}
	return processInput(&in, limit, verbose)
}