package processor

import (
	"crypto/sha512"
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
		atom.P:     1.0,
		atom.Li:    1.0,
		atom.Code:  0.7,
		atom.A:     0.4,
	}
	tagOrder = []atom.Atom{
		atom.Title,
		atom.H1,
		atom.H2,
		atom.H3,
		atom.H4,
		atom.H5,
		atom.H6,
		atom.P,
		atom.Li,
		atom.Code,
		atom.A,
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
// Returns a slice of tags as 1st result,
// a title of the page as 2nd and
// a version of the document based on the hashed contents as 3rd.
//
func ParseHTML(reader io.ReadCloser, verbose, noStopWords bool) ([]*Tag, string, []byte) {
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
		return []*Tag{}, "", nil
	}

	tags, title := collectTags(contents, verbose, noStopWords)

	return tags, title, contents.hash()
}

func crawl(reader io.Reader) *contents {
	contents := &contents{c: make(map[atom.Atom][]string), len: 0}
	crawler := &crawler{}

	z := html.NewTokenizer(reader)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return contents
		case tt == html.StartTagToken:
			token := z.Token()
			if _, ok := tagWeights[token.DataAtom]; ok {
				crawler.push(token.DataAtom)
			}
		case tt == html.EndTagToken:
			token := z.Token()
			if _, ok := tagWeights[token.DataAtom]; ok {
				crawler.pop()
			}
		case tt == html.TextToken:
			if crawler.isCrawling() {
				token := z.Token()
				current := crawler.current()
				if _, ok := contents.c[current]; !ok {
					contents.c[current] = make([]string, 0)
				}
				contents.c[current] = append(contents.c[current], strings.TrimSpace(token.Data))
				contents.len++
			}
		}
	}
}

func collectTags(contents *contents, verbose, noStopWords bool) ([]*Tag, string) {
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
			if tag == atom.Title {
				pageTitle = l
			}
			if isHeading(tag) && l == pageTitle {
				// avoid doubling of scores for duplicated page's title in headings
				if verbose {
					fmt.Printf("<%s>: skipped equal to <title>\n", tag.String())
				}
				continue
			}
			sentences := SplitToSentences(l)
			for _, s := range sentences {
				docsCount++
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

	// Assure page title
	if pageTitle == "" {
		h1s, ok := contents.c[atom.H1]
		if ok && len(h1s) == 1 {
			pageTitle = h1s[0]
		}
	}

	return flatten(tokenIndex), pageTitle
}

func isHeading(t atom.Atom) bool {
	switch t {
	case atom.H1, atom.H2, atom.H3:
		return true
	default:
		return false
	}
}

// contents stores text from target tags.
type contents struct {
	len int
	c   map[atom.Atom][]string
}

func (cnt *contents) forEach(it func(i int, tag atom.Atom, lines []string)) {
	for i, tag := range tagOrder {
		lines, ok := cnt.c[tag]
		if !ok {
			continue
		}
		it(i, tag, lines)
	}
}

func (cnt *contents) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	cnt.forEach(func(i int, tag atom.Atom, lines []string) {
		sb.WriteString(" ")
		sb.WriteString(tag.String())
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%v", lines))
	})
	sb.WriteString(" ]")
	return sb.String()
}

func (cnt *contents) hash() []byte {
	h := sha512.New()
	cnt.forEach(func(i int, tag atom.Atom, lines []string) {
		_, _ = h.Write([]byte(fmt.Sprintf("%s:%v", tag.String(), lines)))
	})
	return h.Sum(nil)
}

// crawler keeps track of the current state of the HTML parser.
type crawler struct {
	stack []atom.Atom
}

func (c *crawler) push(a atom.Atom) {
	c.stack = append(c.stack, a)
}

func (c *crawler) current() atom.Atom {
	if len(c.stack) == 0 {
		return 0
	}
	return c.stack[len(c.stack)-1]
}

func (c *crawler) pop() atom.Atom {
	if len(c.stack) == 0 {
		return 0
	}
	last := len(c.stack) - 1
	v := c.stack[last]
	c.stack = c.stack[:last]
	return v
}

func (c *crawler) isCrawling() bool {
	return len(c.stack) > 0
}
