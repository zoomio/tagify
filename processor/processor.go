package processor

import (
	"math"
	"regexp"
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/zoomio/stopwords"
)

var (
	sanitizeRegex         = regexp.MustCompile(`([^a-z-']*)([a-z-']+)([^a-z-']*)`)
	notAWord              = regexp.MustCompile(`([^a-z'-]+)`)
	doubleNotWordySymbols = regexp.MustCompile(`[\W]{2}`)
)

// sanitize ...
func sanitize(strs []string, noStopWords bool) []string {
	result := make([]string, 0)
	for _, s := range strs {
		normilized, ok := Normalize(s, noStopWords)
		if !ok {
			continue
		}
		result = append(result, normilized)
	}
	return result
}

// Run - 1st sorts given list,
// then iterates over it and de-dupes items in the list by merging inflections,
// then sorts de-duped list again and
// takes only requested size (limit) or just everything if result is smaller than limit.
//
// nolint: gocyclo
func Run(items []*Tag, limit int) []*Tag {
	uniqueTags := make([]*Tag, 0)
	seenTagValues := make(map[string]int)
	uniqueTagsMap := make(map[string]int)

	sortTagItems(items)

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

		// if either item has different singular form and singular form has been seen already OR
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

	sortTagItems(uniqueTags)

	// take only requested size (limit) or just everything if result is smaller than limit
	return uniqueTags[:int(math.Min(float64(limit), float64(len(uniqueTags))))]
}

// Normalize sanitizes word and tells whether it is allowed token or not.
func Normalize(word string, noStopWords bool) (string, bool) {
	// All letters to lower and with proper quote
	word = strings.Replace(strings.ToLower(word), "’", "'", -1)

	// False if it is a stop word
	if noStopWords && stopwords.IsStopWord(word) {
		return word, false
	}

	// False if doesn't match allowed regex
	if !sanitizeRegex.MatchString(word) {
		return word, false
	}

	// Remove not allowed symbols (sanitize)
	word = sanitizeRegex.ReplaceAllString(word, "${2}")

	// Defensive check if sanitized result is still not a word
	if notAWord.MatchString(word) || doubleNotWordySymbols.MatchString(word) {
		return word, false
	}

	// Defensive check if word starts with hyphen
	if strings.HasPrefix(word, "-") {
		return word, false
	}

	// Allowed word
	return word, true
}
