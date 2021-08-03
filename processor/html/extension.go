package html

import "github.com/zoomio/tagify/config"

// HTMLExtension ...
type HTMLExtension interface {
	config.Extension
	Process(cfg *config.Config, line *HTMLLine) error
}
