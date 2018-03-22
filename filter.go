package tagify

import (
	"regexp"
	"strings"

	"github.com/zoomio/tagify/inout"
)

var (
	index stopWords
	reg = regexp.MustCompile(`([^a-z-']*)([a-z-']+)([^a-z-']*)`)
)

// stopWords ...
type stopWords map[string]bool

func indexStopWords(strs []string) stopWords {
	sw := stopWords(make(map[string]bool))
	for _, s := range strs {
		sw[strings.ToLower(s)] = true
	}
	return sw
}

// sanitize ...
func sanitize(s string) string {
	if !reg.MatchString(s) {
		return ""
	}
	return reg.ReplaceAllString(s, "${2}")
}

// InitStopWords ...
func InitStopWords() {
	in := inout.NewIn("stop-word-list.txt")
	strs := in.ReadAllStrings()
	index = indexStopWords(strs)
}

// Filter ...
func Filter(strs []string) []string {
	filtered := make([]string, 0)
	for _, s := range strs {
		v := sanitize(strings.ToLower(s))
		if v == "" || index[v] {
			continue
		}
		filtered = append(filtered, v)
	}
	return filtered
}