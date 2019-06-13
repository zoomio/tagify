package tagify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	tags, err := GetTags("http://stackoverflow.com", HTML, 10, false, false)
	assert.Nil(t, err)
	assert.Len(t, tags, 10)
}

func TestGetTagsFromString(t *testing.T) {
	tags, err := GetTagsFromString("Test input reader of type text", Text, 3, false, true)
	assert.Nil(t, err)
	assert.Len(t, tags, 3)
}

func TestToStrings(t *testing.T) {
	tags, _ := GetTagsFromString("Test input reader of type text", Text, 3, false, true)
	strs := ToStrings(tags)
	assert.Len(t, strs, 3)
}
