package html

import (
	"fmt"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/extension"
)

// HTMLExtension ...
type HTMLExtension interface {
	extension.Extension
	Run(cfg *config.Config, line *HTMLLine) error
}

// HTMLExtensions ...
func HTMLExtensions(exts []extension.Extension) []HTMLExtension {
	res := []HTMLExtension{}
	for _, v := range exts {
		if e, ok := v.(HTMLExtension); ok {
			res = append(res, e)
		}
	}
	return res
}

func RunExtensions(cfg *config.Config, line *HTMLLine, exts []HTMLExtension) {
	for _, v := range exts {
		if cfg.Verbose {
			fmt.Printf("running %q %s\n", v.Name(), v.Version())
		}
		err := v.Run(cfg, line)
		if err != nil {
			fmt.Printf("error in running %q %s: %v\n", v.Name(), v.Version(), err)
		}
	}
}
