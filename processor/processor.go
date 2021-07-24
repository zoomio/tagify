package processor

import (
	"math"

	"github.com/jinzhu/inflection"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/processor/model"
	"github.com/zoomio/tagify/processor/util"
)

// Run - 1st sorts given list,
// then iterates over it and de-dupes items in the list by merging inflections,
// then sorts de-duped list again and
// takes only requested size (limit) or just everything if result is smaller than limit.
//
// nolint: gocyclo
func Run(c *config.Config, items []*model.Tag) []*model.Tag {
	uniqueTags := make([]*model.Tag, 0)
	seenTagValues := make(map[string]int)
	uniqueTagsMap := make(map[string]int)

	util.SortTagItems(items)

	for i, tag := range items {

		// collect indexes of seen items
		if _, ok := seenTagValues[tag.Value]; !ok {
			seenTagValues[tag.Value] = i
		}

		var singularForm string
		if c.Lang == "en" || c.Lang == "" {
			singularForm = inflection.Singular(tag.Value)
		} else {
			singularForm = tag.Value
		}

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
			uniqueTags[savedIndex] = &model.Tag{
				Value:     saved.Value,
				Score:     saved.Score + tag.Score,
				Count:     saved.Count + tag.Count,
				Docs:      saved.Docs + tag.Docs,
				DocsCount: saved.DocsCount,
			}
		}
	}

	// apply TF-IDF
	for _, t := range uniqueTags {
		if t.Docs > 0 && t.DocsCount > 0 {
			t.Score = util.TFIDF(t)
		}
	}

	util.SortTagItems(uniqueTags)

	// take only requested size (limit) or just everything if result is smaller than limit
	result := uniqueTags[:int(math.Min(float64(c.Limit), float64(len(uniqueTags))))]

	// adjust scores to the interval of 0.0 to 1.0
	if c.AdjustScores {
		maxScore := uniqueTags[0].Score
		for _, t := range result {
			t.Score = t.Score / maxScore
		}
	}

	return result
}
