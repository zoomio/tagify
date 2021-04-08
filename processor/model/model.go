package model

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// InputReader ...
type InputReader interface {
	ReadLines() ([]string, error)
	io.ReadCloser
}

// Tag holds some arbitrary string value (e.g. a word) along with some extra data about it.
type Tag struct {
	// Value of the tag, i.e. a word
	Value string
	// Score used to represent importance of the tag
	Score float64
	// Count is the number of times tag appeared in a text
	Count int
	// Docs is the number of documents in a text in which the tag appeared
	Docs int
	// DocsCount is the number of documents in a text
	DocsCount int
}

// Wight input types
const (
	String TagWeightsType = iota // <tagName1>:<tagScore1>|<tagName2>:<tagScore2>
	JSON                         // { "<tagName1>": <tagScore1>, "<tagName2>": <tagScore2> }
)

// TagWeightsType ...
type TagWeightsType byte

// TagWeights ...
type TagWeights map[string]float64

func ParseTagWeights(reader io.Reader, readerType TagWeightsType) TagWeights {
	weights := TagWeights{}

	switch readerType {
	case String:
		buf := new(strings.Builder)
		if _, err := io.Copy(buf, reader); err != nil {
			println(fmt.Errorf("error: can't read string: %w", err))
		}
		for _, v := range strings.Split(buf.String(), "|") {
			tuple := strings.Split(v, ":")
			if len(tuple) != 2 {
				continue
			}
			f, err := strconv.ParseFloat(tuple[1], 64)
			if err != nil {
				println(fmt.Errorf("error: can't read score for [%s]: %w", tuple[0], err))
			}
			weights[tuple[0]] = f
		}
	case JSON:
		if err := json.NewDecoder(reader).Decode(&weights); err != nil {
			println(fmt.Errorf("error: can't read JSON: %w\n", err))
		}
	default:
		fmt.Printf("error: unknown readerType\n")
	}

	return weights
}

func (t *Tag) String() string {
	return fmt.Sprintf("(%s - [score: %.2f, count: %d, docs: %d, docs_count: %d])",
		t.Value, t.Score, t.Count, t.Docs, t.DocsCount)
}

// ParseOutput is a result of the `ParseFunc`.
type ParseOutput struct {
	Tags     map[string]*Tag
	DocTitle string
	DocHash  []byte
	Err      error
}

// FlatTags transforms internal token register into a slice.
func (po *ParseOutput) FlatTags() []*Tag {
	return flatten(po.Tags)
}

type ParseConfig struct {
	Verbose     bool
	NoStopWords bool
	ContentOnly bool
	FullSite    bool
	Source      string
	TagWeights
}

// ParseFunc represents an arbitrary handler,
// which goes through given reader and produces tags.
type ParseFunc func(reader io.ReadCloser, options ...ParseOption) *ParseOutput

func flatten(dict map[string]*Tag) []*Tag {
	flat := make([]*Tag, len(dict))
	var i int
	for _, val := range dict {
		flat[i] = val
		i++
	}
	return flat
}

// ToStrings transforms list of given tags into a list of strings.
func ToStrings(items []*Tag) []string {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = item.Value
	}
	return strs
}
