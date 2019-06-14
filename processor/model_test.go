package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	tag := &Tag{
		Value: "foo",
		Score: 2.5,
		Count: 3,
	}
	assert.Equal(t, "(foo - [score: 2.50, count: 3])", tag.String())
}

func Test_flatten(t *testing.T) {
	dict := map[string]*Tag{
		"foo": {
			Value: "foo",
			Score: 2.5,
			Count: 3,
		},
		"bar": {
			Value: "bar",
			Score: 1.5,
			Count: 2,
		},
	}

	tags := flatten(dict)

	assert.Len(t, tags, 2)
	assert.Equal(t, "foo", tags[0].Value)
	assert.Equal(t, "bar", tags[1].Value)
}

func Test_ToStrings(t *testing.T) {
	dict := []*Tag{
		{
			Value: "foo",
			Score: 2.5,
			Count: 3,
		},
		{
			Value: "bar",
			Score: 1.5,
			Count: 2,
		},
	}

	strs := ToStrings(dict)

	assert.Len(t, strs, 2)
	assert.Equal(t, "foo", strs[0])
	assert.Equal(t, "bar", strs[1])
}
