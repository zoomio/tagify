package rank

import (
	"math"
	"regexp"
	"strings"

	"github.com/jinzhu/inflection"
)

var (
	index = stopWords(make(map[string]bool))
	reg   = regexp.MustCompile(`([^a-z-']*)([a-z-']+)([^a-z-']*)`)
)

// stopWords ...
type stopWords map[string]bool

// sanitize ...
func sanitize(s string) string {
	if !reg.MatchString(s) {
		return ""
	}
	return reg.ReplaceAllString(s, "${2}")
}

func toSingular(s string) string {
	return inflection.Singular(s)
}

// RegisterStopWords ...
func RegisterStopWords(words []string) {
	for _, s := range words {
		index[strings.ToLower(s)] = true
	}
}

// Filter ...
func Filter(strs []string) []string {
	filtered := make([]string, 0)
	for _, s := range strs {
		sanitized, ok := Normalize(s)
		if !ok {
			continue
		}
		filtered = append(filtered, sanitized)
	}
	return filtered
}

// Normalize sanitizes word and tells whether it is allowed token or not.
func Normalize(sanitized string) (string, bool) {
	v := sanitize(strings.ToLower(sanitized))
	if v == "" {
		return v, false
	}
	if index[v] {
		return v, false
	}
	return v, true
}

// Dedupe ...
func Dedupe(items []*Item, limit int) []string {
	var deDupedSize int
	var selectedItems []*Item
	lenItems := len(items)
	if limit <= 0 || limit >= lenItems {
		deDupedSize = lenItems
		selectedItems = items
	} else {
		deDupedSize = int(math.Min(float64(lenItems), float64(limit)))
		selectedItems = items[:int(math.Min(float64(lenItems), float64(limit*2)))]
	}
	return deDupe(deDupedSize, selectedItems)
}

func deDupe(size int, items []*Item) []string {
	index := make(map[string]int)
	deDuped := make([]string, size)
	var i, j int
	for i < size {
		item := items[j]
		if _, ok := index[item.Value]; !ok {
			index[item.Value] = j
		}
		s := toSingular(item.Value)
		k, ok := index[s]
		if s != item.Value && (!ok || ok && j < k) {
			deDuped[i] = item.Value
			index[s] = j
			i++
		} else if s == item.Value && j == k {
			deDuped[i] = item.Value
			i++
		}
		j++
	}
	return deDuped
}
