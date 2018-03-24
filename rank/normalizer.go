package rank

import (
	"regexp"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/jinzhu/inflection"

	"github.com/zoomio/tagify/inout"
)

var (
	index stopWords
	reg   = regexp.MustCompile(`([^a-z-']*)([a-z-']+)([^a-z-']*)`)
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

func toSingular(s string) string {
	return inflection.Singular(s)
}

// InitStopWords ...
func InitStopWords(box *packr.Box) {
	in := inout.NewInFromString(box.String("stop-word-list.txt"))
	strs := in.ReadAllStrings()
	index = indexStopWords(strs)
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
	return toSingular(v), true
}
