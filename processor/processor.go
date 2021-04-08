package processor

import (
	"math"

	"github.com/jinzhu/inflection"
	"github.com/zoomio/stopwords"

	"github.com/zoomio/tagify/processor/model"
	"github.com/zoomio/tagify/processor/util"
)

func init() {
	stopwords.Setup(
		stopwords.Words(stopwords.StopWordsRu),
		stopwords.Words(stopwords.StopWordsZh),
		stopwords.Words(stopwords.StopWordsJa),
		stopwords.Words(stopwords.StopWordsKo),
		stopwords.Words(stopwords.StopWordsHi),
		stopwords.Words(stopwords.StopWordsHe),
		stopwords.Words(stopwords.StopWordsAr),
		stopwords.Words(stopwords.StopWordsDe),
		stopwords.Words(stopwords.StopWordsEs),
		stopwords.Words(stopwords.StopWordsFr),
	)
}

// Run - 1st sorts given list,
// then iterates over it and de-dupes items in the list by merging inflections,
// then sorts de-duped list again and
// takes only requested size (limit) or just everything if result is smaller than limit.
//
// nolint: gocyclo
func Run(items []*model.Tag, limit int) []*model.Tag {
	uniqueTags := make([]*model.Tag, 0)
	seenTagValues := make(map[string]int)
	uniqueTagsMap := make(map[string]int)

	util.SortTagItems(items)

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
			uniqueTags[savedIndex] = &model.Tag{
				Value:     saved.Value,
				Score:     saved.Score + tag.Score,
				Count:     saved.Count + tag.Count,
				Docs:      saved.Docs + tag.Docs,
				DocsCount: saved.DocsCount,
			}
		}
	}

	// Apply TF-IDF
	for _, t := range uniqueTags {
		if t.Docs > 0 && t.DocsCount > 0 {
			t.Score = util.TFIDF(t)
		}
	}

	util.SortTagItems(uniqueTags)

	// take only requested size (limit) or just everything if result is smaller than limit
	return uniqueTags[:int(math.Min(float64(limit), float64(len(uniqueTags))))]
}
