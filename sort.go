package tagify

import (
	"sort"
)

// sortByCountDescending sorts items by count in descending order.
func sortByCountDescending(items []*item) {
	by(func(i1, i2 *item) bool {
		if i1.count > i2.count {
			return true
		}
		return false
	}).Sort(items)
}

// ------------------------------------------ Sort ------------------------------------------

// by is the type of a "less" function that defines the ordering of its item arguments.
type by func(i1, i2 *item) bool

// Sort is a method on the function type, by, that sorts the argument slice according to the function.
func (by by) Sort(items []*item) {
	sorter := &itemSorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(sorter)
}

// itemSorter joins a by function and a slice of Items to be sorted.
type itemSorter struct {
	items []*item
	by    func(p1, p2 *item) bool // Closure used in the Less method.
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
