package processor

import (
	"fmt"
	"strings"

	"github.com/gpestana/htmlizer"
)

var (
	tagWeights = map[string]float64{
		"<h1>": 2,
		"<h2>": 1.3,
		"<h3>": 1.2,
		"<h4>": 1.1,
		"<p>":  1,
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
func ParseHTML(lines []string, tagify, verbose, doFiltering bool) ([]string, []*Tag) {
	// will trim out all the tabs from text
	hizer, err := htmlizer.New([]rune{'\t'})
	if err != nil {
		panic(fmt.Sprintf("error in parsing HTML lines: %v", err))
	}

	for _, line := range lines {
		err = hizer.Load(line)
		if err != nil {
			fmt.Printf("error in loading line \"%s\": %v", line, err)
		}
	}

	if verbose {
		fmt.Println("\nparsed HTML: ")
		fmt.Printf("%v\n\n", hizer)
	}

	textLines := make([]string, 0)
	index := make(map[string]*Tag)

	for tag, weight := range tagWeights {
		tags, err := hizer.GetValues(tag)
		if err != nil {
			fmt.Printf("error in getting values of tag %s: %v", tag, err)
			continue
		}
		for _, t := range tags {
			textLines = append(textLines, t.Value)

			if !tagify {
				continue
			}

			tokens := sanitize(strings.Fields(t.Value), doFiltering)
			for _, token := range tokens {
				item, ok := index[token]
				if !ok {
					item = &Tag{Value: token}
					index[token] = item
				}
				item.Score += weight
				item.Count++
			}
		}
	}

	return textLines, flatten(index)
}
