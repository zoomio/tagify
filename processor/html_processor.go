package processor

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	htmlTagWeights = map[atom.Atom]float64{
		atom.Title:  3,
		atom.H1:     2,
		atom.H2:     1.5,
		atom.H3:     1.4,
		atom.H4:     1.3,
		atom.H5:     1.2,
		atom.H6:     1.1,
		atom.P:      1.0,
		atom.B:      1.2,
		atom.U:      1.2,
		atom.Strong: 1.2,
		atom.I:      1.1,
		atom.Li:     1.0,
		atom.Code:   0.7,
		atom.A:      0.6,
	}
	htmlTagOrder = []atom.Atom{
		atom.Title,
		atom.H1,
		atom.H2,
		atom.H3,
		atom.H4,
		atom.H5,
		atom.H6,
		atom.P,
		atom.B,
		atom.U,
		atom.Strong,
		atom.I,
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
var ParseHTML ParseFunc = func(reader io.ReadCloser, options ...ParseOption) *ParseOutput {

	c := &parseConfig{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	if c.verbose {
		fmt.Println("--> parsing HTML...")
	}

	defer reader.Close()
	contents := parseHTML(reader)

	if c.verbose {
		fmt.Println("--> parsed")
		fmt.Printf("%s\n", contents)
	}

	if len(contents.lines) == 0 {
		return &ParseOutput{}
	}

	tags, title := tagifyHTML(contents, c.verbose, c.noStopWords)

	return &ParseOutput{Tags: tags, DocTitle: title, DocHash: contents.hash()}
}

func parseHTML(reader io.Reader) *htmlContents {
	contents := &htmlContents{lines: make([]*htmlLine, 0)}
	parser := &htmlParser{}

	z := html.NewTokenizer(reader)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// end of the document, we're done
			return contents
		case tt == html.StartTagToken:
			token := z.Token()
			if _, ok := htmlTagWeights[token.DataAtom]; ok {
				parser.push(token.DataAtom)
			}
		case tt == html.EndTagToken:
			token := z.Token()
			if _, ok := htmlTagWeights[token.DataAtom]; ok {
				parser.pop()
				if isHTMLContent(token.DataAtom) || parser.isEmpty() {
					parser.lineIndex++
				}
			}
		case tt == html.TextToken:
			if parser.isNotEmpty() {
				token := z.Token()
				data := bytes.TrimSpace([]byte(token.Data))
				cur := parser.current()

				// skip empty or unknown lines
				if len(data) == 0 || cur == 0 {
					continue
				}

				contents.append(parser.lineIndex, cur, data)
			}
		}
	}
}

func tagifyHTML(contents *htmlContents, verbose, noStopWords bool) ([]*Tag, string) {
	tokenIndex := make(map[string]*Tag)
	var docsCount int
	var pageTitle string

	if verbose {
		fmt.Println("--> tokenized")
	}

	for _, l := range contents.lines {
		s := string(l.data())

		if l.tag == atom.Title {
			pageTitle = s
		} else if l.tag == atom.H1 && pageTitle == "" {
			pageTitle = s
		} else if isHTMLHeading(l.tag) && s == pageTitle {
			// avoid doubling of scores for duplicated page's title in headings
			if verbose {
				fmt.Printf("<%s>: skipped equal to <title>\n", l.tag.String())
			}
			continue
		}

		sentences := l.sentences()
		for _, s := range sentences {
			// skip random non-text related tags
			if !isHTMLContent(s.tag) {
				continue
			}

			docsCount++
			visited := map[string]bool{}

			s.forEach(func(i int, p *htmlPart) {
				weight := htmlTagWeights[p.tag]
				tokens := sanitize(bytes.Fields(p.data), noStopWords)
				if verbose && len(tokens) > 0 {
					fmt.Printf("<%s>: %v\n", l.tag.String(), tokens)
				}

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
			})

			// increment number of appearances in documents for each visited tag
			for token := range visited {
				tokenIndex[token].Docs++
			}
		}
	}

	// set total number of dicuments in the text.
	for _, v := range tokenIndex {
		v.DocsCount = docsCount
	}

	return flatten(tokenIndex), pageTitle
}

func isHTMLHeading(t atom.Atom) bool {
	switch t {
	case atom.H1, atom.H2, atom.H3:
		return true
	default:
		return false
	}
}

func isHTMLContent(t atom.Atom) bool {
	switch t {
	case atom.Title, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.P:
		return true
	default:
		return false
	}
}

// htmlContents stores text from target tags.
type htmlContents struct {
	lines []*htmlLine
}

func (cnt *htmlContents) append(lineIndex int, k atom.Atom, v []byte) {
	for len(cnt.lines) <= lineIndex {
		cnt.lines = append(cnt.lines, &htmlLine{tag: k, parts: make([]*htmlPart, 0)})
	}
	line := cnt.lines[lineIndex]
	line.parts = append(line.parts, &htmlPart{tag: k, data: v})
}

func (cnt *htmlContents) forEach(it func(i int, line *htmlLine)) {
	for i, l := range cnt.lines {
		// skip unsupported tags
		if _, ok := htmlTagWeights[l.tag]; !ok {
			continue
		}
		it(i, l)
	}
}

func (cnt *htmlContents) String() string {
	var sb strings.Builder
	cnt.forEach(func(i int, line *htmlLine) {
		sb.WriteString(fmt.Sprintf("[%d] ", i))
		sb.WriteString(line.String())
		sb.WriteString("\n")
	})
	return sb.String()
}

func (cnt *htmlContents) hash() []byte {
	h := sha512.New()
	cnt.forEach(func(i int, line *htmlLine) {
		_, _ = h.Write([]byte(line.tag.String()))
		_, _ = h.Write([]byte(":"))
		line.forEach(func(i int, p *htmlPart) {
			_, _ = h.Write([]byte(p.tag.String()))
			_, _ = h.Write([]byte(":"))
			_, _ = h.Write(p.data)
		})
	})
	return h.Sum(nil)
}

type htmlPart struct {
	tag  atom.Atom
	data []byte
}

func (d *htmlPart) String() string {
	return fmt.Sprintf("<%s>: %s", d.tag.String(), string(d.data))
}

type htmlLine struct {
	tag   atom.Atom
	parts []*htmlPart
}

func (l *htmlLine) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<%s> - %d parts: [ ", l.tag.String(), len(l.parts)))
	l.forEach(func(i int, p *htmlPart) {
		sb.WriteString("'")
		sb.WriteString(p.String())
		sb.WriteString("' ")
	})
	sb.WriteString("]")
	return sb.String()
}

func (l *htmlLine) forEach(it func(i int, p *htmlPart)) {
	for i, p := range l.parts {
		it(i, p)
	}
}

func (l *htmlLine) data() []byte {
	bs := []byte{}
	l.forEach(func(i int, p *htmlPart) {
		bs = append(bs, p.data...)
	})
	return bs
}

func (l *htmlLine) sentences() []*htmlLine {
	ret := []*htmlLine{}
	data := l.data()
	split := punctuationRegex.ReplaceAll(bytes.TrimSpace(data), newLine)
	sents := bytes.Split(split, newLine)
	var sentArea, partsOffset, pi int
	for _, s := range sents {
		sentArea += len(s)
		snt := &htmlLine{tag: l.tag, parts: []*htmlPart{}}
		ret = append(ret, snt)
		for pi < len(l.parts) {
			partSize := len(l.parts[pi].data)
			if partsOffset+partSize > sentArea {
				break
			}
			snt.parts = append(snt.parts, l.parts[pi])
			partsOffset += partSize
			pi++
		}
	}
	return ret
}

// htmlParser keeps track of the current state of the HTML parser.
type htmlParser struct {
	lineIndex int
	stack     []atom.Atom
}

func (p *htmlParser) head() atom.Atom {
	if len(p.stack) == 0 {
		return 0
	}
	return p.stack[0]
}

func (p *htmlParser) current() atom.Atom {
	if len(p.stack) == 0 {
		return 0
	}
	return p.stack[len(p.stack)-1]
}

func (p *htmlParser) push(a atom.Atom) {
	p.stack = append(p.stack, a)
}

func (p *htmlParser) pop() atom.Atom {
	if len(p.stack) == 0 {
		return 0
	}
	last := len(p.stack) - 1
	v := p.stack[last]
	p.stack = p.stack[:last]
	return v
}

func (p *htmlParser) isEmpty() bool {
	return len(p.stack) == 0
}

func (p *htmlParser) isNotEmpty() bool {
	return len(p.stack) > 0
}
