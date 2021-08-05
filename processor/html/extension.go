package html

import (
	"fmt"

	"golang.org/x/net/html"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
)

// HTMLExtension ...
type HTMLExtension interface {
	extension.Extension
	ParseTag(cfg *config.Config, token *html.Token, lineIdx int) error
	ParseText(cfg *config.Config, token *html.Token, lineIdx int) error
	Tagify(cfg *config.Config, line *HTMLLine) error
}

// HTMLExtensions ...
func extHTML(exts []extension.Extension) []HTMLExtension {
	res := []HTMLExtension{}
	for _, v := range exts {
		if e, ok := v.(HTMLExtension); ok {
			res = append(res, e)
		}
	}
	return res
}

func extParseTag(cfg *config.Config, exts []HTMLExtension, token *html.Token, lineIdx int) {
	for _, v := range exts {
		if cfg.Verbose {
			fmt.Printf("parsing HTML tag %q %s\n", v.Name(), v.Version())
		}
		err := v.ParseTag(cfg, token, lineIdx)
		if err != nil {
			fmt.Printf("error in parsing HTML tag %q %s: %v\n", v.Name(), v.Version(), err)
		}
	}
}

func extParseText(cfg *config.Config, exts []HTMLExtension, token *html.Token, lineIdx int) {
	for _, v := range exts {
		if cfg.Verbose {
			fmt.Printf("parsing HTML text %q %s\n", v.Name(), v.Version())
		}
		err := v.ParseText(cfg, token, lineIdx)
		if err != nil {
			fmt.Printf("error in parsing HTML text %q %s: %v\n", v.Name(), v.Version(), err)
		}
	}
}

func extTagify(cfg *config.Config, exts []HTMLExtension, line *HTMLLine) {
	for _, v := range exts {
		if cfg.Verbose {
			fmt.Printf("tagifying %q %s\n", v.Name(), v.Version())
		}
		err := v.Tagify(cfg, line)
		if err != nil {
			fmt.Printf("error in tagifying %q %s: %v\n", v.Name(), v.Version(), err)
		}
	}
}
