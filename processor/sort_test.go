package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortTagItems_byScore(t *testing.T) {
	tags := []*Tag{
		{Value: "foo", Score: 0.1, Count: 1},
		{Value: "bar", Score: 0.8, Count: 1},
		{Value: "bee", Score: 0.4, Count: 1},
	}

	sortTagItems(tags)

	assert.Equal(t, "bar", tags[0].Value)
	assert.Equal(t, "bee", tags[1].Value)
	assert.Equal(t, "foo", tags[2].Value)
}

func Test_sortTagItems_byCountIfScoreIsEqual(t *testing.T) {
	tags := []*Tag{
		{Value: "foo", Score: 0.1, Count: 1},
		{Value: "bar", Score: 0.1, Count: 3},
		{Value: "bee", Score: 0.1, Count: 2},
	}

	sortTagItems(tags)

	assert.Equal(t, "bar", tags[0].Value)
	assert.Equal(t, "bee", tags[1].Value)
	assert.Equal(t, "foo", tags[2].Value)
}

func Test_sortTagItems_byValueIfScoreAndCountAreEqual(t *testing.T) {
	tags := []*Tag{
		{Value: "foo", Score: 0.1, Count: 1},
		{Value: "bar", Score: 0.1, Count: 1},
		{Value: "bee", Score: 0.1, Count: 1},
	}

	sortTagItems(tags)

	assert.Equal(t, "bar", tags[0].Value)
	assert.Equal(t, "bee", tags[1].Value)
	assert.Equal(t, "foo", tags[2].Value)
}
