package html

import (
	"fmt"

	"golang.org/x/net/html"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
)

// HTMLExt ...
type HTMLExt interface {
	extension.Extension
}

type HTMLExtParseTag interface {
	HTMLExt

	// ParseTag returns true in case if the contnets have been appended and false otherwise.
	ParseTag(cfg *config.Config, token *html.Token, lineIdx int, cnts *HTMLContents) (bool, error)
}

type HTMLExtParseText interface {
	HTMLExt

	// ParseText ...
	ParseText(cfg *config.Config, tagName, text string, lineIdx int) error
}

type HTMLExtTagify interface {
	HTMLExt
	Tagify(cfg *config.Config, line *HTMLLine) error
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

func extParseTag(cfg *config.Config, exts []HTMLExt, token *html.Token, lineIdx int, cnts *HTMLContents) bool {
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
			fmt.Printf("error in parsing HTML tag %q %s: %v\n", v.Name(), v.Version(), err)
		}
		if !appended && ok {
			appended = true
		}
	}
	return appended
}

func extParseText(cfg *config.Config, exts []HTMLExt, tagName, text string, lineIdx int) {
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
			fmt.Printf("error in parsing HTML text %q %s: %v\n", v.Name(), v.Version(), err)
		}
	}
}

func extTagify(cfg *config.Config, exts []HTMLExt, line *HTMLLine) {
	for _, v := range exts {
		e, ok := v.(HTMLExtTagify)
		if !ok {
			continue
		}
		if cfg.Verbose {
			fmt.Printf("tagifying %q %s\n", v.Name(), v.Version())
		}
		err := e.Tagify(cfg, line)
		if err != nil {
			fmt.Printf("error in tagifying %q %s: %v\n", v.Name(), v.Version(), err)
		}
	}
}
