package processor

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	tagWeights = map[atom.Atom]float64{
		atom.H1: 2,
		atom.H2: 1.5,
		atom.H3: 1.4,
		atom.H4: 1.3,
		atom.H5: 1.2,
		atom.H6: 1.1,
		atom.P:  0.9,
		atom.A:  1,
	}
)

type contents struct {
	len int
	c   map[atom.Atom][]string
}

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
func ParseHTML(reader io.Reader, verbose, noStopWords bool) []*Tag {
	contents := crawl(reader)

	if verbose {
		fmt.Println("parsed: ")
		fmt.Printf("%v\n", contents.c)
	}

	if contents.len == 0 {
		return []*Tag{}
	}

	return collectTags(contents, verbose, noStopWords)
}

func crawl(reader io.Reader) *contents {
	contents := &contents{c: make(map[atom.Atom][]string), len: 0}

	z := html.NewTokenizer(reader)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return contents
		case tt == html.StartTagToken:
			t := z.Token()

			if _, ok := tagWeights[t.DataAtom]; ok {
				tt := z.Next()

				if tt == html.TextToken {
					next := z.Token()
					if _, ok := contents.c[t.DataAtom]; !ok {
						contents.c[t.DataAtom] = make([]string, 0)
					}
					contents.c[t.DataAtom] = append(contents.c[t.DataAtom], strings.TrimSpace(next.Data))
					contents.len++
				}
			}
		}
	}
}

func collectTags(contents *contents, verbose, noStopWords bool) []*Tag {
	tagIndex := make(map[string]*Tag)

	for tag, weight := range tagWeights {
		lines, ok := contents.c[tag]
		if !ok {
			continue
		}
		if verbose && lines != nil && len(lines) > 0 {
			fmt.Printf("reading tag: %s\n", tag.String())
		}
		for _, l := range lines {
			tokens := sanitize(strings.Fields(l), noStopWords)
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
