package extension

// Extension provides ability to extend Tagify workflow,
// to be able to incorporate more functionality into the Tagify
// and build such things like deeper more opinionated Tagify primitives
// e.g. YouTube extension to get more data out of the YouTube videos and etc
// (see testImgCrawlerExt inside the processor/html/extension_test.go).
type Extension interface {
	Name() string
	Version() string
	Result() *ExtResult
}

type BaseExtension struct {
	name    string
	version string
	*ExtResult
}

func NewExtension(name, version string) *BaseExtension {
	return &BaseExtension{name: name, version: version}
}

func (ext *BaseExtension) Name() string {
	return ext.name
}

func (ext *BaseExtension) Version() string {
	return ext.version
}

func (ext *BaseExtension) Result() *ExtResult {
	return ext.ExtResult
}

// ExtResult ...
type ExtResult struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Err     error                  `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// NewResult ...
func NewResult(ext Extension, data map[string]interface{}, err error) *ExtResult {
	return &ExtResult{
		Name:    ext.Name(),
		Version: ext.Version(),
		Err:     err,
		Data:    data,
	}
}

// MapResults ...
func MapResults(exts []Extension) map[string]map[string]*ExtResult {
	res := map[string]map[string]*ExtResult{}
	for _, v := range exts {
		e, ok := res[v.Name()]
		if !ok {
			e = map[string]*ExtResult{}
			res[v.Name()] = e
		}
		e[v.Version()] = v.Result()
	}
	return res
}
