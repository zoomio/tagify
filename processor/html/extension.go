package html

import (
	"fmt"

	"golang.org/x/net/html"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
	"github.com/zoomio/tagify/model"
)

const HTMLParseEndErrorMsg = "received stop command, exiting HTML parser"

// HTMLExt ...
type HTMLExt interface {
	extension.Extension
}

// HTMLExtParseTag executed at the HTML parsing phase when dealing with the HTML tag.
type HTMLExtParseTag interface {
	HTMLExt

	// ParseTag returns true in case if the contents have been appended and false otherwise.
	ParseTag(cfg *config.Config, token *html.Token, lineIdx int, cnts *HTMLContents) (bool, error)
}

// HTMLExtParseText executed at the HTML parsing phase when dealing with the text inside an HTML tag.
type HTMLExtParseText interface {
	HTMLExt

	// ParseText ...
	ParseText(cfg *config.Config, tagName, text string, lineIdx int) error
}

// HTMLExtParseText executed during token counting phase.
type HTMLExtTagify interface {
	HTMLExt
	Tagify(cfg *config.Config, line *HTMLLine, tokenIndex map[string]*model.Tag) error
}

func NewHTMLParseEndError() *HTMLParseEndError {
	return &HTMLParseEndError{}
}

type HTMLParseEndError struct {
}

func (e *HTMLParseEndError) Error() string {
	return HTMLParseEndErrorMsg
}

// HTMLExtensions ...
func extHTML(exts []extension.Extension) []HTMLExt {
	res := []HTMLExt{}
	for _, v := range exts {
		if e, ok := v.(HTMLExt); ok {
			res = append(res, e)
		}
	}
	return res
}

func extParseTag(cfg *config.Config, exts []HTMLExt, token *html.Token, lineIdx int, cnts *HTMLContents) (bool, error) {
	var appended bool
	for _, v := range exts {
		e, ok := v.(HTMLExtParseTag)
		if !ok {
			continue
		}
		if cfg.Verbose {
			fmt.Printf("parsing HTML tag %q %s\n", v.Name(), v.Version())
		}
		ok, err := e.ParseTag(cfg, token, lineIdx, cnts)
		if err != nil {
			if cfg.Verbose {
				fmt.Printf("error in parsing HTML tag %q in %q %s: %v\n", token.DataAtom.String(), v.Name(), v.Version(), err)
			}
			return appended, err
		}
		if !appended && ok {
			appended = true
		}
	}
	return appended, nil
}

func extParseText(cfg *config.Config, exts []HTMLExt, tagName, text string, lineIdx int) error {
	for _, v := range exts {
		e, ok := v.(HTMLExtParseText)
		if !ok {
			continue
		}
		if cfg.Verbose {
			fmt.Printf("parsing HTML text %q %s\n", v.Name(), v.Version())
		}
		err := e.ParseText(cfg, tagName, text, lineIdx)
		if err != nil {
			if cfg.Verbose {
				fmt.Printf("error in parsing HTML text %q %s: %v\n", v.Name(), v.Version(), err)
			}
			return err
		}
	}
	return nil
}

func extTagify(cfg *config.Config, exts []HTMLExt, line *HTMLLine, tokenIndex map[string]*model.Tag) {
	for _, v := range exts {
		e, ok := v.(HTMLExtTagify)
		if !ok {
			continue
		}
		if cfg.Verbose {
			fmt.Printf("tagifying %q %s\n", v.Name(), v.Version())
		}
		err := e.Tagify(cfg, line, tokenIndex)
		if err != nil && cfg.Verbose {
			fmt.Printf("error in tagifying %q %s: %v\n", v.Name(), v.Version(), err)
		}
	}
}
