package config

import (
	"github.com/zoomio/tagify/processor/html"
	"github.com/zoomio/tagify/processor/model"
)

// Extension ...
type Extension interface {
	Name() string
	Version() string
	Result() *ExtResult
}

// HTMLExtension ...
type HTMLExtension interface {
	Extension
	Process(cfg *Config, line *html.HTMLLine) error
}

// ExtResult ...
type ExtResult struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}
