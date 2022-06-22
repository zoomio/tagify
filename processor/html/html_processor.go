package html

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/zoomio/stopwords"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
	"github.com/zoomio/tagify/model"
	"github.com/zoomio/tagify/processor/util"
)

var (
	// default weights for HTML tags
	defaultTagWeights = config.TagWeights{
		"h1":     2,
		"title":  1.7,
		"meta":   1.7,
		"h2":     1.5,
		"h3":     1.4,
		"h4":     1.3,
		"h5":     1.2,
		"h6":     1.1,
		"p":      1.0,
		"b":      1.2,
		"u":      1.2,
		"strong": 1.2,
		"i":      1.1,
		"li":     1.0,
		"code":   0.7,
		"a":      0.6,
	}
)

func isHTMLHeading(t string) bool {
	switch atom.Lookup([]byte(t)) {
	case atom.H1, atom.H2, atom.H3:
		return true
	default:
		return false
	}
}

func isHTMLContent(t string) bool {
	switch atom.Lookup([]byte(t)) {
	case atom.Title, atom.Meta, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.P:
		return true
	default:
		return false
	}
}

func isNonClosingSingleTag(t string) bool {
	switch atom.Lookup([]byte(t)) {
	case atom.Meta, atom.Link:
		return true
	default:
		return false
	}
}

func isSameDomain(href, domain string) bool {
	if strings.HasPrefix(href, "/") {
		return true
	}

	// double trim http(s)
	dest := strings.TrimPrefix(strings.TrimPrefix(href, "https://"), "http://")
	// double trim http(s)
	host := strings.TrimPrefix(strings.TrimPrefix(domain, "https://"), "http://")
	// first slash in the address
	i := strings.Index(dest, "/")

	if i == -1 {
		return strings.HasSuffix(dest, host)
	}

	return strings.HasSuffix(dest[:i], host)
}

func updateDetectStr(candidate, controlStr string) string {
	if len(candidate) > len(controlStr) {
		return candidate
	}
	return controlStr
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
var ProcessHTML model.ProcessFunc = func(c *config.Config, reader io.ReadCloser) *model.Result {

	defer reader.Close()

	if c.Verbose {
		fmt.Println("--> parsing HTML...")
	}

	var err error
	var contents *HTMLContents
	var parseFn parseFunc = ParseHTML

	exts := extHTML(c.Extensions)

	if c.TagWeights == nil {
		c.TagWeights = defaultTagWeights
	}
	if c.ExtraTagWeights != nil {
		for k, v := range c.ExtraTagWeights {
			c.TagWeights[k] = v
		}
	}

	if c.Verbose {
		fmt.Printf("using configuration: %#v\n", c)
	}

	if c.FullSite && c.Source != "" {
		var crawler *webCrawler
		crawler, err = newWebCrawler(parseFn, exts, c.Source, c.Verbose)
		if err != nil {
			return model.ErrResult(err)
		}
		contents = crawler.run(reader)
	} else {
		contents = parseFn(reader, c, exts, nil)
	}

	if c.Verbose {
		fmt.Println("--> parsed")
	}

	if err != nil {
		return &model.Result{Err: err}
	}

	if len(contents.lines) == 0 {
		return model.EmptyResult()
	}

	tags, title, lang := tagifyHTML(contents, c, exts)

	return &model.Result{
		Meta: &model.Meta{
			ContentType: config.HTML,
			DocTitle:    title,
			DocHash:     fmt.Sprintf("%x", contents.hash()),
			Lang:        lang,
		},
		RawTags:    tags,
		Extensions: extension.MapResults(c.Extensions),
	}
}

func ParseHTML(reader io.Reader, cfg *config.Config, exts []HTMLExt, c *webCrawler) *HTMLContents {
	contents := &HTMLContents{lines: make([]*HTMLLine, 0), htmlTagWeights: cfg.TagWeights}
	parser := &htmlParser{}

	var controlStr string

	if !cfg.SkipLang {
		defer func() {
			// detect language and setup stop words for it
			if cfg.StopWords == nil {
				info := whatlanggo.Detect(controlStr)
				contents.lang = info.Lang.String()
				cfg.SetStopWords(info.Lang.Iso6391())
				if cfg.Verbose {
					fmt.Printf("detected language based on %q: %s [%s] [%s]\n ",
						controlStr, info.Lang.String(), info.Lang.Iso6391(), info.Lang.Iso6393())
				}
				if cfg.NoStopWords {
					contents.reg = cfg.StopWords
				}
			}
		}()
	}

	var cursor string
	z := html.NewTokenizer(reader)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// end of the document, we're done
			return contents
		case html.SelfClosingTagToken: // e.g. <img ... />
			if parser.shouldStop() {
				// flag has been set to true, exiting
				if cfg.Verbose {
					fmt.Println(HTMLParseEndErrorMsg)
				}
				return contents
			}
			token := z.Token()
			if _, ok := cfg.TagWeights[token.Data]; ok || cfg.AllTagWeights {
				_, err := extParseTag(cfg, exts, &token, parser.lineIndex, contents)
				if err != nil {
					switch err.(type) {
					case *HTMLParseEndError:
						if cfg.Verbose {
							fmt.Println(err.Error())
						}
						return contents
					}
				}
			}
		case html.StartTagToken:
			token := z.Token()
			cursor = token.Data
			if _, ok := cfg.TagWeights[cursor]; !ok && !cfg.AllTagWeights {
				continue
			}

			parser.push(cursor)

			var appended bool

			// handle <meta name="description" content="...">
			if cursor == atom.Meta.String() {
				var name, content string
				for _, a := range token.Attr {
					if a.Key == "name" {
						name = a.Val
					}
					if a.Key == "content" {
						content = a.Val
					}
					if name != "" && content != "" {
						break
					}
				}
				if name == "description" {
					controlStr = updateDetectStr(content, controlStr)
					contents.Append(parser.lineIndex, cursor, []byte(content))
					appended = true
				}
			}

			// allow for extensions
			ok, err := extParseTag(cfg, exts, &token, parser.lineIndex, contents)
			if err != nil {
				switch err.(type) {
				case *HTMLParseEndError:
					parser.stop()
				}
			}
			if !appended && ok {
				appended = true
			}

			// handle non-self-closing nor has closing tags, e.g. <link ...> & <meta ...>
			if isNonClosingSingleTag(cursor) {
				parser.pop()
				if appended {
					parser.lineIndex++
				}

				if parser.shouldStop() {
					// flag has been set to true, exiting
					if cfg.Verbose {
						fmt.Println(HTMLParseEndErrorMsg)
					}
					return contents
				}
			}

			// go follow links in case if web crawler is ON.
			if c != nil && cursor == atom.A.String() {
				for _, a := range token.Attr {
					if a.Key == "href" && isSameDomain(a.Val, c.domain) {
						c.crawl(a.Val)
						break
					}
				}
			}

		case html.EndTagToken:
			token := z.Token()
			if _, ok := cfg.TagWeights[token.Data]; ok || cfg.AllTagWeights {
				parser.pop()
				cursor = parser.current()
				if parser.isEmpty() {
					parser.lineIndex++
				}
			}

			if parser.shouldStop() {
				// flag has been set to true, exiting
				if cfg.Verbose {
					fmt.Println(HTMLParseEndErrorMsg)
				}
				return contents
			}
		case html.TextToken:
			token := z.Token()
			if _, ok := cfg.TagWeights[cursor]; (ok || cfg.AllTagWeights) && parser.isNotEmpty() {

				// skip empty or unknown lines
				if len(strings.TrimSpace(token.Data)) == 0 {
					continue
				}

				// Take only <title> from <head> and ignore the rest <title> tags in the body
				if cursor == atom.Title.String() && parser.parent() != "" {
					continue
				}

				err := extParseText(cfg, exts, cursor, token.Data, parser.lineIndex)
				if err != nil {
					switch err.(type) {
					case *HTMLParseEndError:
						parser.stop()
					}
				}

				controlStr = updateDetectStr(token.Data, controlStr)
				contents.Append(parser.lineIndex, cursor, []byte(token.Data))
			}
		}
	}
}

func tagifyHTML(contents *HTMLContents, cfg *config.Config,
	exts []HTMLExt) (tokenIndex map[string]*model.Tag, pageTitle string, lang string) {

	tokenIndex = map[string]*model.Tag{}
	lang = contents.lang
	reg := contents.reg

	var docsCount int

	for _, l := range contents.lines {
		s := string(l.data)

		// Title tags have special treatment
		// if (l.tag == atom.Title.String() || l.tag == atom.H1.String()) && len(pageTitle) < len(s) {
		if (l.tag == atom.Title.String() || l.tag == atom.H1.String()) && pageTitle == "" {
			pageTitle = s
		}

		// avoid doubling of scores for duplicated page's title in headings
		if isHTMLHeading(l.tag) && s == pageTitle {
			if cfg.Verbose {
				fmt.Printf("<%s>: skipped equal to <title>\n", l.tag)
			}
			continue
		}

		sentences := l.sentences()
		for _, snt := range sentences {
			// skip random non-text related tags
			if cfg.ContentOnly && !isHTMLContent(snt.tag) {
				continue
			}

			if len(snt.data) == 0 {
				continue
			}

			docsCount++
			visited := map[string]bool{}

			snt.forEach(func(i int, p *htmlPart) {
				var weight float64
				if l.weightOverride {
					weight = l.weight
				} else {
					weight = cfg.TagWeights[p.tag]
				}
				tokens := util.Sanitize(bytes.Fields(snt.pData(p)), reg)

				for _, token := range tokens {
					visited[token] = true
					item, ok := tokenIndex[token]
					if !ok {
						item = &model.Tag{Value: token}
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

		// run extensions if any
		extTagify(cfg, exts, l, tokenIndex)
	}

	// set total number of documents in the text.
	for _, v := range tokenIndex {
		v.DocsCount = docsCount
	}

	return
}

// HTMLContents stores text from target tags.
type HTMLContents struct {
	lines          []*HTMLLine
	htmlTagWeights config.TagWeights

	lang string
	reg  *stopwords.Register
}

func (cnt *HTMLContents) Len() int {
	return len(cnt.lines)
}

func (cnt *HTMLContents) Append(lineIndex int, tag string, data []byte) {
	for len(cnt.lines) <= lineIndex {
		cnt.lines = append(cnt.lines, &HTMLLine{tag: tag, parts: make([]*htmlPart, 0)})
	}
	line := cnt.lines[lineIndex]
	line.add(tag, data)
}

func (cnt *HTMLContents) Weigh(lineIndex int, weight float64) {
	if lineIndex >= len(cnt.lines) {
		// silently ignore for now...
		return
	}
	line := cnt.lines[lineIndex]
	line.weightOverride = true
	line.weight = weight
}

func (cnt *HTMLContents) forEach(it func(i int, line *HTMLLine)) {
	for i, l := range cnt.lines {
		// skip unsupported tags
		if _, ok := cnt.htmlTagWeights[l.tag]; !ok {
			continue
		}
		it(i, l)
	}
}

func (cnt *HTMLContents) String() string {
	var sb strings.Builder
	cnt.forEach(func(i int, line *HTMLLine) {
		sb.WriteString(fmt.Sprintf("[%d] ", i))
		sb.WriteString(line.String())
		sb.WriteString("\n")
	})
	return sb.String()
}

func (cnt *HTMLContents) hash() []byte {
	h := sha512.New()
	cnt.forEach(func(i int, line *HTMLLine) {
		_, _ = h.Write([]byte(line.tag))
		_, _ = h.Write([]byte(":"))
		line.forEach(func(i int, p *htmlPart) {
			_, _ = h.Write([]byte(p.tag))
			_, _ = h.Write([]byte(":"))
			_, _ = h.Write(line.pData(p))
		})
	})
	return h.Sum(nil)
}

// htmlPart is a part of an HTML tag text.
type htmlPart struct {
	tag string
	pos int
	len int
}

func (d *htmlPart) String() string {
	return fmt.Sprintf("<%s>: pos - %d, len - %d", d.tag, d.pos, d.len)
}

type HTMLLine struct {
	tag            string
	parts          []*htmlPart
	data           []byte
	weightOverride bool
	weight         float64
}

func (l *HTMLLine) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<%s> - %d parts: [ ", l.tag, len(l.parts)))
	l.forEach(func(i int, p *htmlPart) {
		sb.WriteString("'")
		sb.WriteString("<")
		sb.WriteString(p.tag)
		sb.WriteString(">: ")
		sb.WriteString(string(l.pData(p)))
		sb.WriteString("' ")
	})
	sb.WriteString("]")
	return sb.String()
}

func (l *HTMLLine) forEach(it func(i int, p *htmlPart)) {
	for i, p := range l.parts {
		it(i, p)
	}
}

func (l *HTMLLine) add(tag string, data []byte) {
	l.parts = append(l.parts, &htmlPart{tag: tag, pos: len(l.data), len: len(data)})
	l.data = append(l.data, data...)
}

func (l *HTMLLine) pData(part *htmlPart) []byte {
	return l.data[part.pos : part.pos+part.len]
}

// breaksdown an HTML line into a slice of HTML sentences.
func (l *HTMLLine) sentences() []*HTMLLine {
	ret := []*HTMLLine{}
	var offset, diff, pDiff, i, j int
	sents := util.SplitToSentences(l.data)
	for i < len(l.parts) && j < len(sents) {
		s := &HTMLLine{tag: l.tag, parts: []*htmlPart{}}
		ret = append(ret, s)

		sent := sents[j]
		sentSize := len(sent)

		part := l.parts[i]
		partSize := part.len

		diff = (offset + partSize) - (offset + sentSize)

		if diff > 0 {
			// part is bigger than sentence, splitting part
			s.add(part.tag, sent)
			offset += sentSize
			pDiff = diff
			j++ // increment index for the next sentence
		} else if diff < 0 {
			// sentence is bigger than part, appending part included into sentence
			if pDiff > 0 {
				s.add(part.tag, l.data[part.pos+pDiff:part.pos+part.len])
				offset += (partSize - pDiff)
				pDiff = 0
			} else {
				s.add(part.tag, l.pData(part))
				offset += partSize
			}
			i++ // increment index for the next part
		} else {
			// part is equal to sentence
			s.add(part.tag, l.pData(part))
			offset += partSize
			pDiff = 0
			i++
			j++
		}
	}
	return ret
}

type htmlParserNode struct {
	name string
	stop bool
}

// htmlParser keeps track of the current state of the HTML parser.
type htmlParser struct {
	lineIndex int
	stack     []*htmlParserNode
}

func (p *htmlParser) current() string {
	if len(p.stack) == 0 {
		return ""
	}
	return p.stack[len(p.stack)-1].name
}

func (p *htmlParser) currentNode() *htmlParserNode {
	if len(p.stack) == 0 {
		return nil
	}
	return p.stack[len(p.stack)-1]
}

func (p *htmlParser) parent() string {
	if len(p.stack) < 2 {
		return ""
	}
	return p.stack[len(p.stack)-2].name
}

func (p *htmlParser) push(a string) {
	p.stack = append(p.stack, &htmlParserNode{name: a})
}

func (p *htmlParser) stop() {
	p.currentNode().stop = true
}

func (p *htmlParser) pop() string {
	if len(p.stack) == 0 {
		return ""
	}
	last := len(p.stack) - 1
	v := p.stack[last]
	p.stack = p.stack[:last]
	return v.name
}

func (p *htmlParser) isEmpty() bool {
	return len(p.stack) == 0
}

func (p *htmlParser) isNotEmpty() bool {
	return len(p.stack) > 0
}

func (p *htmlParser) shouldStop() bool {
	node := p.currentNode()
	if node == nil {
		return false
	}
	return node.stop
}
