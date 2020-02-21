package processor

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// md types
const (
	heading1 mdType = iota
	heading2
	heading3
	heading4
	heading5
	heading6
	paragraph
	boldItalic
	bold
	underscore
	italic
	blockquote
	code
	anchor
	strikethrough
)

var (
	mdTypes = [...]string{
		"heading1",
		"heading2",
		"heading3",
		"heading4",
		"heading5",
		"heading6",
		"paragraph",
		"boldItalic",
		"bold",
		"underscore",
		"italic",
		"blockquote",
		"code",
		"anchor",
		"strikethrough",
	}

	mdWeights = map[mdType]float64{
		heading1:      2,
		heading2:      1.5,
		heading3:      1.4,
		heading4:      1.3,
		heading5:      1.2,
		heading6:      1.1,
		paragraph:     1.0,
		boldItalic:    1.1,
		bold:          1.1,
		underscore:    1.1,
		italic:        1.0,
		blockquote:    1.0,
		code:          0.7,
		anchor:        0.4,
		strikethrough: 0.0,
	}

	boldItalicReg = regexp.MustCompile(`\*\*\*(.*?)\*\*\*`)
	boldReg       = regexp.MustCompile(`\*\*(.*?)\*\*`)
	italicReg     = regexp.MustCompile(`\*(.*?)\*`)
	strikeReg     = regexp.MustCompile(`\~\~(.*?)\~\~`)
	underscoreReg = regexp.MustCompile(`__(.*?)__`)
	anchorReg     = regexp.MustCompile(`\[(.*?)\]\((.*?)\)[^\)]`)
	escapeReg     = regexp.MustCompile(`^\>(\s|)`)
	blockquoteReg = regexp.MustCompile(`\&gt\;(.*?)$`)
	backtipReg    = regexp.MustCompile("`(.*?)`")

	h1Reg = regexp.MustCompile(`^#(\s|)(.*?)$`)
	h2Reg = regexp.MustCompile(`^##(\s|)(.*?)$`)
	h3Reg = regexp.MustCompile(`^###(\s|)(.*?)$`)
	h4Reg = regexp.MustCompile(`^####(\s|)(.*?)$`)
	h5Reg = regexp.MustCompile(`^#####(\s|)(.*?)$`)
	h6Reg = regexp.MustCompile(`^######(\s|)(.*?)$`)
)

type mdType byte

func (t mdType) String() string {
	if t < heading1 || t > strikethrough {
		return "unknown"
	}
	return mdTypes[t]
}

// ParseMD parses given Markdown document input into a slice of tags.
var ParseMD ParseFunc = func(in io.ReadCloser, options ...ParseOption) *ParseOutput {

	c := &parseConfig{}

	// apply custom configuration
	for _, option := range options {
		option(c)
	}

	if c.verbose {
		fmt.Println("--> parsing Markdown...")
	}

	defer in.Close()
	contents := parseMD(in)

	if c.verbose {
		fmt.Println("--> parsed")
		fmt.Printf("%s\n", contents)
	}

	tags, title := tagifyMD(contents, c.verbose, c.noStopWords)

	return &ParseOutput{Tags: tags, DocTitle: title, DocHash: contents.hash()}
}

func parseMD(reader io.Reader) *mdContents {

	contents := &mdContents{lines: make([]*mdLine, 0)}
	scanner := bufio.NewScanner(reader)

	i := -1

	for scanner.Scan() {

		line := bytes.TrimSpace(scanner.Bytes())

		// skip empty lines
		if len(line) == 0 {
			continue
		}

		i++

		// 1. handle headings
		if line[0] == '#' {
			count := bytes.Count(line, []byte(`#`))
			switch count {
			case 1:
				contents.append(i, heading1, h1Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 2:
				contents.append(i, heading2, h1Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 3:
				contents.append(i, heading3, h1Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 4:
				contents.append(i, heading4, h1Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 5:
				contents.append(i, heading5, h1Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 6:
				contents.append(i, heading6, h1Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			}
		}

		// 2. handle quote
		// escape and wrap blockquotes in "<blockquote>" tags
		line = escapeReg.ReplaceAll(line, []byte(`&gt;`))
		if blockquoteReg.Match(line) {
			contents.append(i, blockquote, blockquoteReg.ReplaceAll(line, []byte(`$1`)))
			continue
		}

		// 3. handle paragraph
		handlers := map[int]*mdPartHandler{}
		// wrap bold and italic text in "<b>" and "<i>" elements
		// line = appendLine(i, line, boldItalicReg, paragraph, boldItalic, contents)
		appendHandlers(handlers, line, boldItalic, boldItalicReg)
		// line = appendLine(i, line, boldReg, paragraph, bold, contents)
		appendHandlers(handlers, line, bold, boldReg)
		// line = appendLine(i, line, italicReg, paragraph, italic, contents)
		appendHandlers(handlers, line, italic, italicReg)
		// wrap strikethrough text in "<s>" tags
		// line = appendLine(i, line, strikeReg, paragraph, strikethrough, contents)
		appendHandlers(handlers, line, strikethrough, strikeReg)
		// wrap underscored text in "<u>" tags
		// line = appendLine(i, line, underscoreReg, paragraph, underscore, contents)
		appendHandlers(handlers, line, underscore, underscoreReg)
		// convert links to anchor tags
		// line = appendLine(i, line, anchorReg, paragraph, anchor, contents)
		appendHandlers(handlers, line, anchor, anchorReg)
		// wrap the content of backticks inside of "<code>" tags
		// line = appendLine(i, line, backtipReg, paragraph, code, contents)
		appendHandlers(handlers, line, code, backtipReg)
		println("enter the matrix")
		appendLine(i, line, handlers, contents)
		println("exit the matrix")
	}

	return contents
}

func tagifyMD(contents *mdContents, verbose, noStopWords bool) ([]*Tag, string) {
	tokenIndex := make(map[string]*Tag)
	var docsCount int
	var pageTitle string

	if verbose {
		fmt.Println("--> tokenized")
	}

	for _, l := range contents.lines {
		s := string(l.data())

		if isMDHeading(l.tag) && pageTitle == "" {
			pageTitle = s
		}

		sentences := l.sentences()
		for _, s := range sentences {
			docsCount++
			visited := map[string]bool{}

			s.forEach(func(i int, p *mdPart) {
				weight := mdWeights[p.tag]
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

func isMDHeading(t mdType) bool {
	switch t {
	case heading1, heading2, heading3, heading4, heading5, heading6:
		return true
	default:
		return false
	}
}

func appendLine(lineIndex int, line []byte, handlers map[int]*mdPartHandler, cnt *mdContents) {
	i := 0
	p := []byte{}
	for i < len(line) {
		var ok bool
		var h *mdPartHandler

		if h, ok = handlers[i]; !ok {
			p = append(p, line[i])
			i++
			continue
		}

		if len(p) > 0 {
			cnt.append(lineIndex, paragraph, p)
			p = []byte{}
		}

		cnt.append(lineIndex, h.tag, h.re.ReplaceAll(line[h.start:h.end], []byte(`$1`)))
		i = h.end
	}
}

func appendHandlers(handlers map[int]*mdPartHandler, line []byte, tag mdType, re *regexp.Regexp) {
	pairs := re.FindAllIndex(line, -1)
	for _, pair := range pairs {
		handlers[pair[0]] = &mdPartHandler{start: pair[0], end: pair[1], tag: boldItalic, re: boldItalicReg}
	}
}

type mdPartHandler struct {
	start int
	end   int
	tag   mdType
	re    *regexp.Regexp
}

// mdContents stores text from target tags.
type mdContents struct {
	lines []*mdLine
}

func (cnt *mdContents) append(lineIndex int, tag mdType, data []byte) {
	if len(cnt.lines) <= lineIndex {
		cnt.lines = append(cnt.lines, &mdLine{tag: tag, parts: make([]*mdPart, 0)})
	}

	// skip empty data
	if len(data) == 0 {
		return
	}

	line := cnt.lines[lineIndex]
	line.parts = append(line.parts, &mdPart{tag: tag, data: data})
}

func (cnt *mdContents) forEach(it func(i int, line *mdLine)) {
	for k, v := range cnt.lines {
		it(k, v)
	}
}

func (cnt *mdContents) String() string {
	var sb strings.Builder
	cnt.forEach(func(i int, line *mdLine) {
		sb.WriteString(fmt.Sprintf("[%d] ", i))
		sb.WriteString(line.String())
		sb.WriteString("\n")
	})
	return sb.String()
}

func (cnt *mdContents) hash() []byte {
	h := sha512.New()
	cnt.forEach(func(i int, line *mdLine) {
		_, _ = h.Write(line.bytes())
	})
	return h.Sum(nil)
}

type mdPart struct {
	tag  mdType
	data []byte
}

func (p *mdPart) String() string {
	return fmt.Sprintf("<%s>: %s", p.tag.String(), string(p.data))
}

func (p *mdPart) bytes() []byte {
	return append(p.data, byte(p.tag))
}

type mdLine struct {
	tag   mdType
	parts []*mdPart
}

func (l *mdLine) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<%s> - %d parts: [ ", l.tag.String(), len(l.parts)))
	l.forEach(func(i int, p *mdPart) {
		sb.WriteString("'")
		sb.WriteString(p.String())
		sb.WriteString("' ")
	})
	sb.WriteString("]")
	return sb.String()
}

func (l *mdLine) forEach(it func(i int, p *mdPart)) {
	for i, p := range l.parts {
		it(i, p)
	}
}

func (l *mdLine) bytes() []byte {
	bs := []byte{byte(l.tag)}
	for _, elm := range l.parts {
		bs = append(bs, elm.bytes()...)
	}
	return bs
}

func (l *mdLine) data() []byte {
	bs := []byte{}
	for _, p := range l.parts {
		bs = append(bs, p.data...)
	}
	return bs
}

func (l *mdLine) sentences() []*mdLine {
	ret := []*mdLine{}
	data := l.data()
	split := punctuationRegex.ReplaceAll(bytes.TrimSpace(data), newLine)
	sents := bytes.Split(split, newLine)
	var sentArea, partsOffset, pi int
	for _, s := range sents {
		sentArea += len(s)
		snt := &mdLine{tag: l.tag, parts: []*mdPart{}}
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
