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

// GetTags produces slice of tags ordered by frequency and limited by limit.
func GetTags(source string, limit int, verbose bool) []string {
	in := inout.NewIn(source)

	var items []*rank.Item

	switch in.SourceType {
	case inout.STDIn, inout.FS:
		items = rank.ParseText(in.ReadAllStrings())
	case inout.Web:
		items = rank.ParseHTML(in.GetLines(), verbose)
	default:
		panic("unrecognized source")
	}

	sortByScoreDescending(items)
	if limit > 0 {
		return toStrings(items[:limit])
	}
	return toStrings(items)
}
