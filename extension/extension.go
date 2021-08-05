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
	Data    map[string]interface{} `json:"data,omitempty"`
}

// NewResult ...
func NewResult(ext Extension, data map[string]interface{}, err error) *Result {
	return &Result{
		Name:    ext.Name(),
		Version: ext.Version(),
		Err:     err,
		Data:    data,
	}
}

// GetResults ...
func GetResults(exts []Extension) []*Result {
	res := make([]*Result, len(exts))
	for k, v := range exts {
		res[k] = v.Result()
	}
	return res
}
