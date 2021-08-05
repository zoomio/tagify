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
	ParseTag(cfg *config.Config, token *html.Token, lineIdx int) error
}

type HTMLExtParseText interface {
	HTMLExt
	ParseText(cfg *config.Config, token *html.Token, lineIdx int) error
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

func extParseTag(cfg *config.Config, exts []HTMLExt, token *html.Token, lineIdx int) {
	for _, v := range exts {
		e, ok := v.(HTMLExtParseTag)
		if !ok {
			continue
		}
		if cfg.Verbose {
			fmt.Printf("parsing HTML tag %q %s\n", v.Name(), v.Version())
		}
		err := e.ParseTag(cfg, token, lineIdx)
		if err != nil {
			fmt.Printf("error in parsing HTML tag %q %s: %v\n", v.Name(), v.Version(), err)
		}
	}
}

func extParseText(cfg *config.Config, exts []HTMLExt, token *html.Token, lineIdx int) {
	for _, v := range exts {
		e, ok := v.(HTMLExtParseText)
		if !ok {
			continue
		}
		if cfg.Verbose {
			fmt.Printf("parsing HTML text %q %s\n", v.Name(), v.Version())
		}
		err := e.ParseText(cfg, token, lineIdx)
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
