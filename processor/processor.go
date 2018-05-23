package processor

import (
	"math"
	"regexp"
	"strings"

	"github.com/jinzhu/inflection"
)

const (
	maxLimit     = 20
	defaultLimit = 5
)

var (
	stopWordsIndex = make(map[string]bool)
	sanitizeRegex  = regexp.MustCompile(`([^a-z-']*)([a-z-']+)([^a-z-']*)`)
)

// RegisterStopWords ...
func RegisterStopWords(words []string) {
	for _, s := range words {
		stopWordsIndex[strings.ToLower(s)] = true
	}
}

// Filter ...
func Filter(strs []string) []string {
	result := make([]string, 0)
	for _, s := range strs {
		normilized, ok := normalize(s)
		if !ok {
			continue
		}
		result = append(result, normilized)
	}
	return result
}

// Run first sorts given list based on scores,
// then iterates over the given list and de-dupes items in the list by merging inflections,
// then sorts de-duped list by scores in descending order and
// takes only rquested size (limit) or just everything if result is smaller than limit.
//
// nolint: gocyclo
func Run(items []*Tag, limit int) []*Tag {
	uniqueTags := make([]*Tag, 0)
	seenTagValues := make(map[string]int)
	uniqueTagsMap := make(map[string]int)

	sortByScoreDescending(items)

	for i, tag := range items {

		// collect indexes of seen items
		if _, ok := seenTagValues[tag.Value]; !ok {
			seenTagValues[tag.Value] = i
		}

		singularForm := inflection.Singular(tag.Value)
		seenIndex, seen := seenTagValues[singularForm]

		// if item has different singular form, but singular form hasn't been seen yet,
		// then add current form of item to unique, and set current index for singular form in seenTagValues.
		if tag.Value != singularForm && !seen {
			uniqueTags = append(uniqueTags, tag)
			uniqueTagsMap[singularForm] = len(uniqueTags) - 1
			seenTagValues[singularForm] = i
		}

		// if item has same singular form, and its seen index is the same as curent,
		// then add item to unique.
		if tag.Value == singularForm && seenIndex == i {
			uniqueTags = append(uniqueTags, tag)
			uniqueTagsMap[singularForm] = len(uniqueTags) - 1
		}

		// if either item has different singular form and singular form has been seen already or
		// item is in singular form and has predecessor, then merge scores of both forms into predecessor.
		if (tag.Value != singularForm && seen) || (tag.Value == singularForm && seenIndex < i) {
			savedIndex := uniqueTagsMap[singularForm]
			saved := uniqueTags[savedIndex]
			uniqueTags[savedIndex] = &Tag{
				Value: saved.Value,
				Score: saved.Score + tag.Score,
				Count: saved.Count + tag.Count,
			}
		}
	}

	sortByScoreDescending(uniqueTags)

	if limit == 0 || limit > maxLimit {
		limit = defaultLimit
	}

	// take only rquested size (limit) or just everything if result is smaller than limit
	return uniqueTags[:int(math.Min(float64(limit), float64(len(uniqueTags))))]
}

// normalize sanitizes word and tells whether it is allowed token or not.
func normalize(word string) (string, bool) {
	// All letters to lower
	word = strings.ToLower(word)

	// False if doesn't match allowed regex
	if !sanitizeRegex.MatchString(word) {
		return "", false
	}

	// Remove not allowed symbols (sanitize)
	word = sanitizeRegex.ReplaceAllString(word, "${2}")

	// False if it is a stop word
	if stopWordsIndex[word] {
		return word, false
	}

	// Allowed word
	return word, true
}
