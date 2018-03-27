package rank

import (
	"fmt"
	"strings"

	"github.com/gpestana/htmlizer"
)

var (
	tagWeights = map[string]float64{
		"<h1>": 5,
		"<h2>": 4,
		"<h3>": 3,
		"<h4>": 2,
		"<p>":  1,
	}
)

// ParseHTML receives lines of raw strings from the Web and produces result of prioritised tags
// based on the importance of HTML tags which wrap sentences.
//
// Example:
//	<h1>A story about foo
//	<p> Foo was a good guy but, had a quite poor time management skills,
//	therefore he had issues with shipping all his tasks. Though foo had heaps
//	of other amazing skills, which were appreciated by his management.
//
// Result:
//	foo: 5 + 1, story: 5, management: 1 + 1, skills: 1 + 1.
//
func ParseHTML(lines []string, verbose bool) []*Item {
	// will trim out all the tabs from text
	hizer, err := htmlizer.New([]rune{'\t'})
	if err != nil {
		panic(fmt.Sprintf("error in parsing HTML lines: %v", err))
	}

	for _, line := range lines {
		hizer.Load(line)
	}

	if verbose {
		fmt.Println("\nparsed HTML: ")
		fmt.Printf("%v\n\n", hizer)
	}

	index := make(map[string]*Item)

	for tag, weight := range tagWeights {
		tags, err := hizer.GetValues(tag)
		if err != nil {
			fmt.Printf("error in  getting values of tag %s: %v", tag, err)
			continue
		}
		for _, t := range tags {
			tokens := Filter(strings.Fields(t.Value))
			for _, token := range tokens {
				item, ok := index[token]
				if !ok {
					item = &Item{Value: token}
					index[token] = item
				}
				item.Score = item.Score + weight
			}
		}
	}

	return flatten(index)
}
