package rank

import (
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
	return toSingular(v), true
}
