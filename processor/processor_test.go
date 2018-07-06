package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	n, ok := normalize("(part)")
	assert.True(t, ok)
	assert.Equal(t, "part", n)

	n, ok = normalize("(part1)")
	assert.True(t, ok)
	assert.Equal(t, "part", n)

	n, ok = normalize("part1")
	assert.True(t, ok)
	assert.Equal(t, "part", n)

	n, ok = normalize("part}1")
	assert.True(t, ok)
	assert.Equal(t, "part", n)

	n, ok = normalize("{part}")
	assert.True(t, ok)
	assert.Equal(t, "part", n)

	n, ok = normalize("{d'arko}")
	assert.True(t, ok)
	assert.Equal(t, "d'arko", n)

	n, ok = normalize("1)")
	assert.False(t, ok)
	assert.Equal(t, "", n)
}

func TestSanitize_timestamp(t *testing.T) {
	n, ok := normalize("2018-02-24T12:00:49Z")
	assert.False(t, ok)
	assert.Equal(t, "", n)
}

func TestRun_Limits(t *testing.T) {
	items := []*Tag{
		{Value: "cat", Score: 1},
		{Value: "dog", Score: 1},
		{Value: "foo", Score: 1},
		{Value: "bar", Score: 1},
		{Value: "bee", Score: 1},
	}
	processed := Run(items, 3)
	assert.Len(t, processed, 3)
}

func TestRun_Sorts(t *testing.T) {
	items := []*Tag{
		{Value: "cat", Score: 1},
		{Value: "dog", Score: 2},
		{Value: "foo", Score: 5},
		{Value: "bar", Score: 3},
		{Value: "bee", Score: 4},
	}
	processed := Run(items, 5)
	assert.Len(t, processed, 5)
	assert.Equal(t, "foo", processed[0].Value)
	assert.Equal(t, "bee", processed[1].Value)
	assert.Equal(t, "bar", processed[2].Value)
	assert.Equal(t, "dog", processed[3].Value)
	assert.Equal(t, "cat", processed[4].Value)
}

func TestRun_DeDupes(t *testing.T) {
	items := []*Tag{
		{Value: "cat", Score: 5},
		{Value: "person", Score: 2},
		{Value: "people", Score: 5},
		{Value: "bar", Score: 3},
		{Value: "cats", Score: 1},
	}
	processed := Run(items, 5)
	assert.Len(t, processed, 3)
	assert.Equal(t, "people", processed[0].Value)
	assert.Equal(t, "cat", processed[1].Value)
	assert.Equal(t, "bar", processed[2].Value)
}
