package text

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"io"
	"strings"

	"github.com/zoomio/stopwords"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/model"
	"github.com/zoomio/tagify/processor/util"
)

type TxtContents struct {
	lang string
	reg  *stopwords.Register
}

func (cnt *TxtContents) SetLang(l string) {
	cnt.lang = l
}

func (cnt *TxtContents) SetReg(reg *stopwords.Register) {
	cnt.reg = reg
}

// ProcessText parses given text lines of text into a slice of tags.
var ProcessText model.ProcessFunc = func(c *config.Config, in io.ReadCloser) *model.Result {

	if c.Verbose {
		fmt.Println("parsing plain text...")
	}

	// if c.Verbose {
	// 	fmt.Printf("using configuration: %#v\n", c)
	// }

	var docsCount int

	defer in.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(in)
	inStr := buf.String()
	lines := strings.FieldsFunc(inStr, func(r rune) bool {
		return r == '\n'
	})

	if c.Verbose {
		fmt.Printf("got %d lines\n", len(lines))
	}

	if len(lines) == 0 {
		return &model.Result{}
	}

	contents := &TxtContents{}
	tokenIndex := make(map[string]*model.Tag)
	tokens := make([]string, 0)
	for _, l := range lines {
		// detect language and setup stop words for it
		if !c.SkipLang && c.StopWords == nil && len(l) > 0 {
			config.DetectLang(c, l, contents)
		}
		sentences := util.SplitToSentences([]byte(l))
		for _, s := range sentences {
			docsCount++
			tokens = append(tokens, util.Sanitize(bytes.Fields(s), contents.reg)...)
			visited := map[string]bool{}
			for _, token := range tokens {
				visited[token] = true
				item, ok := tokenIndex[token]
				if !ok {
					item = &model.Tag{Value: token}
					tokenIndex[token] = item
				}
				item.Score++
				item.Count++
			}
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

	return &model.Result{
		RawTags: tokenIndex,
		Meta: &model.Meta{
			ContentType: config.Text,
			DocHash:     fmt.Sprintf("%x", hashTokens(tokens)),
			Lang:        contents.lang,
		},
	}
}

func hashTokens(ts []string) []byte {
	h := sha512.New()
	for _, t := range ts {
		_, _ = h.Write([]byte(t))
	}
	return h.Sum(nil)
}
