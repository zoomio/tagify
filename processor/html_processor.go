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
		atom.Title: 3,
		atom.H1:    2,
		atom.H2:    1.5,
		atom.H3:    1.4,
		atom.H4:    1.3,
		atom.H5:    1.2,
		atom.H6:    1.1,
		atom.A:     1,
		atom.P:     0.9,
		atom.Li:    0.9,
	}
	tagOrder = []atom.Atom{
		atom.Title,
		atom.H1,
		atom.H2,
		atom.H3,
		atom.H4,
		atom.H5,
		atom.H6,
		atom.A,
		atom.P,
		atom.Li,
	}
)

type contents struct {
	len int
	c   map[atom.Atom][]string
}

func (cnt *contents) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for _, tag := range tagOrder {
		lines, ok := cnt.c[tag]
		if !ok {
			continue
		}
		sb.WriteString(" ")
		sb.WriteString(tag.String())
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%v", lines))
	}
	sb.WriteString(" ]")
	return sb.String()
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
func ParseHTML(reader io.ReadCloser, verbose, noStopWords bool) []*Tag {
	if verbose {
		fmt.Println("--> parsing HTML...")
	}

	defer reader.Close()
	contents := crawl(reader)

	if verbose {
		fmt.Println("--> parsed")
		fmt.Printf("%s\n", contents)
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
	tokenIndex := make(map[string]*Tag)
	var docsCount int
	var pageTitle string

	if verbose {
		fmt.Println("--> tokenized")
	}

	for _, tag := range tagOrder {
		weight, ok := tagWeights[tag]
		if !ok {
			continue
		}
		lines, ok := contents.c[tag]
		if !ok {
			continue
		}
		for _, l := range lines {
			sentences := SplitToSentences(l)
			for _, s := range sentences {
				docsCount++
				if isHeading(tag) && s == pageTitle {
					// avoid doubling of scores for duplicated page's title in headings
					if verbose {
						fmt.Printf("<%s>: skipped equal to <title>\n", tag.String())
					}
					continue
				}
				if tag == atom.Title {
					pageTitle = s
				}
				tokens := sanitize(strings.Fields(s), noStopWords)
				if verbose && len(tokens) > 0 {
					fmt.Printf("<%s>: %v\n", tag.String(), tokens)
				}
				visited := map[string]bool{}
				for _, token := range tokens {
					visited[token] = true
					item, ok := tokenIndex[token]
					if !ok {
						item = &Tag{Value: token}
						tokenIndex[token] = item
					}
					item.Score += weight
					item.Count++
				}

				// increment number of appearances in documents for each visited tag
				for token := range visited {
					tokenIndex[token].Docs++
				}
			}
		}
	}

	// set total number of dicuments in the text.
	for _, v := range tokenIndex {
		v.DocsCount = docsCount
	}

	return flatten(tokenIndex)
}

func isHeading(t atom.Atom) bool {
	switch t {
	case atom.H1, atom.H2, atom.H3:
		return true
	default:
		return false
	}
}
