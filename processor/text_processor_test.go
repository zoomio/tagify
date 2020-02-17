package processor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoomio/inout"
)

var (
	text = `
		There was a boy
		Who's name was Jim
	`
	text2 = `
		There was a girl
		Who's name was Anne
	`
)

func Test_ParseText_Empty(t *testing.T) {
	tags, _ := ParseText(inout.NewFromString(""), false, false)
	assert.Len(t, tags, 0)
}

func Test_ParseText_WithStopWords(t *testing.T) {
	tags, _ := ParseText(inout.NewFromString(text), false, false)
	assert.Len(t, tags, 7)
	assert.Subset(t, ToStrings(tags), []string{"there", "was", "a", "boy", "who's", "name", "jim"})
}

func Test_ParseText_NoStopWords(t *testing.T) {
	tags, _ := ParseText(inout.NewFromString(text), false, true)
	assert.Len(t, tags, 2)
	assert.Subset(t, ToStrings(tags), []string{"boy", "jim"})
}

func Test_calculatesVersion(t *testing.T) {
	_, version1 := ParseText(inout.NewFromString(text), false, false)
	assert.NotNil(t, version1)
	assert.Equal(t,
		"323c7bf1fe804151d8c378648061d861554a4ae5d02558ce140c1ee3ff186c37a1600bab87abada56d50148b35c121c8b1abb5db8c13a75e9d676fd9130f3c6a",
		fmt.Sprintf("%x", version1))

	_, version2 := ParseText(inout.NewFromString(text2), false, false)
	assert.NotNil(t, version2)
	assert.Equal(t,
		"2f1ba6d722f14042db22ea7c433d02c8b666b33106b1c18bac00388b0ee3add19f411b234981293864698fc9e9dc9073b966378bcb6f49b1c7f07ca99a17a5cc",
		fmt.Sprintf("%x", version2))
}
