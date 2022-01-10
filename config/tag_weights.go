package config

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
			println(fmt.Errorf("error: can't read JSON: %w", err))
		}
	default:
		fmt.Printf("error: unknown readerType\n")
	}

	return weights
}
