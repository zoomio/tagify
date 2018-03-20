package tagify

import (
	"github.com/zoomio/tagify/inout"
)

type item struct {
	value string
	count int
}

func countStrings(strs []string) map[string]int {
	items := make(map[string]int)
	for _, s := range strs {
		items[s]++
	}
	return items
}

func toItems(strs map[string]int) []*item {
	items := make([]*item, len(strs))
	var i int
	for s, count := range strs {
		items[i] = &item{value: s, count: count}
		i++
	}
	return items
}

func toStrings(items []*item) []string {
	strs := make([]string, len(items))
	var i int
	for _, item := range items {
		strs[i] = item.value
		i++
	}
	return strs
}

// Process ...
func Process(source string, limit int) []string {
	in := inout.NewIn(source)
	strs := in.ReadAllStrings()
	if strs == nil || len(strs) == 0 {
		return []string{}
	}
	items := toItems(countStrings(strs))
	sortByCountDescending(items)
	if limit > 0 {
		return toStrings(items[:limit])
	}
	return toStrings(items)
}
