package extension

// Extension provides ability to extend Tagify workflow,
// to be able to incorporate more functionality into the Tagify
// and build such things like deeper more opinionated Tagify primitives
// e.g. YouTube extension to get more data out of the YouTube videos and etc
// (see testImgCrawlerExt inside the processor/html/extension_test.go).
type Extension interface {
	Name() string
	Version() string
	Result() *Result
}

type BaseExtension struct {
	name    string
	version string
	result  *Result
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

func (ext *BaseExtension) Result() string {
	return ext.Result()
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

// MapResults ...
func MapResults(exts []Extension) map[string]map[string]*Result {
	res := map[string]map[string]*Result{}
	for _, v := range exts {
		e, ok := res[v.Name()]
		if !ok {
			e = map[string]*Result{}
			res[v.Name()] = e
		}
		e[v.Version()] = v.Result()
	}
	return res
}
