package extension

// Extension ...
type Extension interface {
	Name() string
	Version() string
	Result() *Result
}

// ExtResult ...
type Result struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Err     error                  `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data"`
}
