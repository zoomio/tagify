package md

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/zoomio/stopwords"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/model"
	"github.com/zoomio/tagify/processor/util"
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

	// default weights for MD tags
	defaultTagWeights = config.TagWeights{
		"heading1":      2,
		"heading2":      1.5,
		"heading3":      1.4,
		"heading4":      1.3,
		"heading5":      1.2,
		"heading6":      1.1,
		"paragraph":     1.0,
		"boldItalic":    1.1,
		"bold":          1.1,
		"underscore":    1.1,
		"italic":        1.0,
		"blockquote":    1.0,
		"code":          0.7,
		"anchor":        0.4,
		"strikethrough": 0.0,
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

// ProcessMD parses given Markdown document input into a slice of tags.
var ProcessMD model.ProcessFunc = func(c *config.Config, in io.ReadCloser) *model.Result {

	if c.Verbose {
		fmt.Println("--> parsing Markdown...")
	}

	defer in.Close()
	contents := ParseMD(in)

	if c.Verbose {
		fmt.Println("--> parsed")
		fmt.Printf("%s\n", contents)
	}

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

	tags, title, lang := tagifyMD(contents, c)

	return &model.Result{
		RawTags: tags,
		Meta: &model.Meta{
			ContentType: config.Markdown,
			DocTitle:    title,
			DocHash:     fmt.Sprintf("%x", contents.hash()),
			Lang:        lang,
		},
	}
}

func ParseMD(reader io.Reader) *MDContents {

	contents := &MDContents{lines: make([]*mdLine, 0)}
	scanner := bufio.NewScanner(reader)

	index := -1

	for scanner.Scan() {

		line := bytes.TrimSpace(scanner.Bytes())
		// skip empty line
		if len(line) == 0 {
			continue
		}

		index++

		// 1. handle headings
		if line[0] == '#' {
			count := bytes.Count(line, []byte(`#`))
			switch count {
			case 1:
				contents.append(index, heading1, h1Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 2:
				contents.append(index, heading2, h2Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 3:
				contents.append(index, heading3, h3Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 4:
				contents.append(index, heading4, h4Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 5:
				contents.append(index, heading5, h5Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			case 6:
				contents.append(index, heading6, h6Reg.ReplaceAll(line, []byte(`$2`)))
				continue
			}
		}

		// 2. handle quote
		// escape and wrap blockquotes in "<blockquote>" tags
		line = escapeReg.ReplaceAll(line, []byte(`&gt;`))
		if blockquoteReg.Match(line) {
			contents.append(index, blockquote, blockquoteReg.ReplaceAll(line, []byte(`$1`)))
			continue
		}

		// 3. handle paragraph
		handlers := map[int]*mdPartHandler{}
		// wrap bold and italic text in "<b>" and "<i>" elements
		appendHandlers(handlers, line, boldItalic, boldItalicReg)
		appendHandlers(handlers, line, bold, boldReg)
		appendHandlers(handlers, line, italic, italicReg)
		// wrap strikethrough text in "<s>" tags
		appendHandlers(handlers, line, strikethrough, strikeReg)
		// wrap underscored text in "<u>" tags
		appendHandlers(handlers, line, underscore, underscoreReg)
		// convert links to anchor tags
		appendHandlers(handlers, line, anchor, anchorReg)
		// wrap the content of backticks inside of "<code>" tags
		appendHandlers(handlers, line, code, backtipReg)
		appendLine(index, line, handlers, contents)
	}

	return contents
}

func tagifyMD(contents *MDContents, c *config.Config) (tokenIndex map[string]*model.Tag, pageTitle string, lang string) {
	tokenIndex = make(map[string]*model.Tag)
	var docsCount int
	var reg *stopwords.Register

	for _, line := range contents.lines {
		// skip empty lines
		if len(line.parts) == 0 {
			continue
		}

		s := string(line.data)

		// detect language and setup stop words for it
		if c.StopWords == nil && s != "" {
			info := whatlanggo.Detect(s)
			lang = info.Lang.String()
			c.SetStopWords(info.Lang.Iso6391())
			if c.Verbose {
				fmt.Printf("detected language: %s [%s] [%s]\n ",
					info.Lang.String(), info.Lang.Iso6391(), info.Lang.Iso6393())
			}
			if c.NoStopWords {
				reg = c.StopWords
			}
		}

		if isMDHeading(line.tag) && pageTitle == "" {
			pageTitle = s
		}

		sentences := line.sentences()
		for _, snt := range sentences {
			if len(snt.data) == 0 {
				continue
			}

			docsCount++
			visited := map[string]bool{}

			snt.forEach(func(i int, p *mdPart) {
				weight := c.TagWeights[p.tag.String()]
				tokens := util.Sanitize(bytes.Fields(snt.pData(p)), reg)
				if c.Verbose && len(tokens) > 0 {
					fmt.Printf("<%s>: %v\n", line.tag.String(), tokens)
				}

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
	}

	// set total number of dicuments in the text.
	for _, v := range tokenIndex {
		v.DocsCount = docsCount
	}

	return
}

func isMDHeading(t mdType) bool {
	switch t {
	case heading1, heading2, heading3, heading4, heading5, heading6:
		return true
	default:
		return false
	}
}

func appendLine(lineIndex int, line []byte, handlers map[int]*mdPartHandler, cnt *MDContents) {
	i := 0
	p := []byte{}
	for i < len(line) {
		h, ok := handlers[i]
		if !ok {
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

	// take care of simple text line cases
	if len(p) > 0 {
		cnt.append(lineIndex, paragraph, p)
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

// MDContents stores text from target tags.
type MDContents struct {
	lines []*mdLine
}

func (cnt *MDContents) append(index int, tag mdType, data []byte) {
	for len(cnt.lines) <= index {
		cnt.lines = append(cnt.lines, &mdLine{tag: tag, parts: make([]*mdPart, 0)})
	}
	line := cnt.lines[index]
	line.add(tag, data)
}

func (cnt *MDContents) forEach(it func(i int, line *mdLine)) {
	for k, v := range cnt.lines {
		it(k, v)
	}
}

func (cnt *MDContents) String() string {
	var sb strings.Builder
	cnt.forEach(func(i int, line *mdLine) {
		sb.WriteString(fmt.Sprintf("[%d] ", i))
		sb.WriteString(line.String())
		sb.WriteString("\n")
	})
	return sb.String()
}

func (cnt *MDContents) hash() []byte {
	h := sha512.New()
	cnt.forEach(func(i int, line *mdLine) {
		_, _ = h.Write([]byte(line.tag.String()))
		_, _ = h.Write([]byte(":"))
		line.forEach(func(i int, p *mdPart) {
			_, _ = h.Write([]byte(p.tag.String()))
			_, _ = h.Write([]byte(":"))
			_, _ = h.Write(line.pData(p))
		})
	})
	return h.Sum(nil)
}

type mdPart struct {
	tag mdType
	pos int
	len int
}

func (p *mdPart) String() string {
	return fmt.Sprintf("<%s>: pos - %d, len - %d", p.tag.String(), p.pos, p.len)
}

type mdLine struct {
	tag   mdType
	parts []*mdPart
	data  []byte
}

func (l *mdLine) add(tag mdType, data []byte) {
	l.parts = append(l.parts, &mdPart{tag: tag, pos: len(l.data), len: len(data)})
	l.data = append(l.data, data...)
}

func (l *mdLine) pData(part *mdPart) []byte {
	return l.data[part.pos : part.pos+part.len]
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

// breaksdown a markdown line into a slice of markdown sentences.
func (l *mdLine) sentences() []*mdLine {
	ret := []*mdLine{}
	var offset, diff, pDiff, i, j int
	sents := util.SplitToSentences(l.data)
	for i < len(l.parts) && j < len(sents) {
		s := &mdLine{tag: l.tag, parts: []*mdPart{}}
		ret = append(ret, s)

		sent := sents[j]
		sentSize := len(sent)

		part := l.parts[i]
		partSize := part.len

		diff = (offset + partSize) - (offset + sentSize)

		if diff > 0 {
			// MD part is bigger than sentence, splitting MD part
			s.add(part.tag, sent)
			offset += sentSize
			pDiff = diff
			j++ // increment index for the next sentence
		} else if diff < 0 {
			// sentence is bigger than MD part, appending MD part included into sentence
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
			// MD part is equal to sentence
			s.add(part.tag, l.pData(part))
			offset += partSize
			pDiff = 0
			i++
			j++
		}
	}
	return ret
}
