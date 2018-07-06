package processor

import (
	"sort"
	"strings"
)

// sortTagItems sorts items by score in descending order,
// if scores are equal it sorts by count if counts are equal, 
// it sorts string values alphabetically.
func sortTagItems(items []*Tag) {
	by(func(i1, i2 *Tag) bool {
		// Higher score goes 1st
		if i1.Score > i2.Score {
			return true
		} else if i1.Score < i2.Score {
			return false
		}

		// Bigger count goes 1st
		if i1.Count > i2.Count {
			return true
		} else if i1.Count < i2.Count {
			return false
		}

		// Alphabetic sort
		return strings.Compare(i1.Value, i2.Value) < 0
	}).Sort(items)
}

// ------------------------------------------ Sort ------------------------------------------

// by is the type of a "less" function that defines the ordering of its item arguments.
type by func(i1, i2 *Tag) bool

// Sort is a method on the function type, by, that sorts the argument slice according to the function.
func (by by) Sort(items []*Tag) {
	sorter := &itemSorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(sorter)
}

// itemSorter joins a by function and a slice of Items to be sorted.
type itemSorter struct {
	items []*Tag
	by    func(p1, p2 *Tag) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *itemSorter) Len() int {
	return len(s.items)
}

// Swap is part of sort.Interface.
func (s *itemSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *itemSorter) Less(i, j int) bool {
	return s.by(s.items[i], s.items[j])
}
