package processor

import (
	"fmt"
	"strings"

	"github.com/gpestana/htmlizer"
)

var (
	tagWeights = map[string]float64{
		"<h1>": 2,
		"<h2>": 1.5,
		"<h3>": 1.4,
		"<h4>": 1.3,
		"<h5>": 1.2,
		"<h6>": 1.1,
		"<p>":  1,
		"<a>":  1,
	}
)

// ParseHTML receives lines of raw HTML markup text from the Web and returns simple text,
// plus list of prioritised tags (if tagify == true)
// based on the importance of HTML tags which wrap sentences.
//
// Example:
//	<h1>A story about foo
//	<p> Foo was a good guy but, had a quite poor time management skills,
//	therefore he had issues with shipping all his tasks. Though foo had heaps
//	of other amazing skills, which gained him a fortune.
//
// Result:
//	foo: 2 + 1 = 3, story: 2, management: 1 + 1 = 2, skills: 1 + 1 = 2.
//
func ParseHTML(html []string, verbose, noStopWords bool) []*Tag {
	// will trim out all the tabs from text
	hizer, err := htmlizer.New([]rune{'\t'})
	if err != nil && verbose {
		fmt.Printf("error in parsing HTML lines: %v\n", err)
		return []*Tag{}
	}

	for _, line := range html {
		err = hizer.Load(line)
		if err != nil && verbose {
			fmt.Printf("error in loading line \"%s\": %v\n", line, err)
		}
	}

	if verbose {
		fmt.Println("\nparsed HTML: ")
		fmt.Printf("%v\n\n", hizer)
	}

	return collectTags(hizer, verbose, noStopWords)
}

func collectTags(hizer htmlizer.Htmlizer, verbose, noStopWords bool) []*Tag {
	tagIndex := make(map[string]*Tag)

	for tag, weight := range tagWeights {
		tags, err := hizer.GetValues(tag)
		if err != nil && verbose {
			fmt.Printf("error in getting values for tag %s: %v\n", tag, err)
			continue
		}
		for _, t := range tags {
			tokens := sanitize(strings.Fields(t.Value), noStopWords)
			for _, token := range tokens {
				item, ok := tagIndex[token]
				if !ok {
					item = &Tag{Value: token}
					tagIndex[token] = item
				}
				item.Score += weight
				item.Count++
			}
		}
	}

	return flatten(tagIndex)
}
