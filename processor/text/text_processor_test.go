package text

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoomio/inout"

	"github.com/zoomio/tagify/processor/model"
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
	out := ParseText(inout.NewFromString(""))
	assert.Len(t, out.Tags, 0)
}

func Test_ParseText_WithStopWords(t *testing.T) {
	out := ParseText(inout.NewFromString(text))
	assert.Len(t, out.Tags, 7)
	assert.Subset(t, model.ToStrings(out.FlatTags()), []string{"there", "was", "a", "boy", "who's", "name", "jim"})
}

func Test_ParseText_NoStopWords(t *testing.T) {
	out := ParseText(inout.NewFromString(text), model.NoStopWords(true))
	assert.Len(t, out.Tags, 2)
	assert.Subset(t, model.ToStrings(out.FlatTags()), []string{"boy", "jim"})
}

func Test_calculatesVersion(t *testing.T) {
	out1 := ParseText(inout.NewFromString(text))
	assert.Nil(t, out1.Err)
	assert.Equal(t,
		"323c7bf1fe804151d8c378648061d861554a4ae5d02558ce140c1ee3ff186c37a1600bab87abada56d50148b35c121c8b1abb5db8c13a75e9d676fd9130f3c6a",
		fmt.Sprintf("%x", out1.DocHash))

	out2 := ParseText(inout.NewFromString(text2))
	assert.Nil(t, out2.Err)
	assert.Equal(t,
		"2f1ba6d722f14042db22ea7c433d02c8b666b33106b1c18bac00388b0ee3add19f411b234981293864698fc9e9dc9073b966378bcb6f49b1c7f07ca99a17a5cc",
		fmt.Sprintf("%x", out2.DocHash))
}
