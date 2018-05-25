package tagify

import (
	"github.com/gobuffalo/packr"

	"github.com/zoomio/tagify/processor"
)

func processInput(in *In, limit int, verbose bool) ([]*processor.Tag, error) {
	var items []*processor.Tag

	switch in.ContentType {
	case HTML:
		lines, err := in.ReadAllLines()
		if err != nil {
			return items, err
		}
		items = processor.ParseHTML(lines, verbose)
	default:
		strs, err := in.ReadAllStrings()
		if err != nil {
			return items, err
		}
		items = processor.ParseText(strs)
	}
	return processor.Run(items, limit), nil
}

// Init initializes Tagify.
func Init() error {
	box := packr.NewBox("./_resources")
	in := NewInFromString(box.String("stop-word-list.txt"), Text)
	strs, err := in.ReadAllStrings()
	if err != nil {
		return err
	}
	processor.RegisterStopWords(strs)
	return nil
}

// GetTags produces slice of tags ordered by frequency and limited by limit.
func GetTags(source string, contentType ContentType, limit int, verbose bool) ([]*processor.Tag, error) {
	in, err := NewIn(source)
	if err != nil {
		return []*processor.Tag{}, err
	}
	if contentType > Unknown {
		in.ContentType = contentType
	}
	return processInput(&in, limit, verbose)
}

// GetTagsFromString produces slice of tags ordered by frequency and limited by limit.
func GetTagsFromString(input string, contentType ContentType, limit int, verbose bool) ([]*processor.Tag, error) {
	in := NewInFromString(input, contentType)
	return processInput(&in, limit, verbose)
}

// ToStrings ...
func ToStrings(items []*processor.Tag) []string {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = item.Value
	}
	return strs
}
