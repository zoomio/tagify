package config

// Extension ...
type Extension interface {
	Name() string
	Version() string
	Result() *ExtResult
}

// ExtResult ...
type ExtResult struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}
