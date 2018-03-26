package tagify

import (
	"fmt"
	
	"github.com/gobuffalo/packr"

	"github.com/zoomio/tagify/rank"
)

func processInput(in *In, limit int, verbose bool) ([]string, error) {
	var items []*rank.Item

	switch in.ContentType {
	case HTML:
		items = rank.ParseHTML(in.GetLines(), verbose)		
	default:
		items = rank.ParseText(in.ReadAllStrings())
	}

	sortByScoreDescending(items)
	
	return rank.Dedupe(items, limit), nil
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
func GetTags(source string, contentType ContentType, limit int, verbose bool) ([]string, error) {
	in, err := NewIn(source)
	if err != nil {
		return []string{}, err
	}
	if contentType > Unknown {
		in.ContentType = contentType
	}
	return processInput(&in, limit, verbose)
}

// GetTagsFromString produces slice of tags ordered by frequency and limited by limit.
func GetTagsFromString(input string, contentType ContentType, limit int, verbose bool) ([]string, error) {
	in, err := NewInFromString(input, contentType)
	if err != nil {
		return []string{}, err
	}
	return processInput(&in, limit, verbose)
}