package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zoomio/tagify/processor/model"
)

func Test_Run_Limits(t *testing.T) {
	items := []*model.Tag{
		{Value: "cat", Score: 1},
		{Value: "dog", Score: 1},
		{Value: "foo", Score: 1},
		{Value: "bar", Score: 1},
		{Value: "bee", Score: 1},
	}
	processed := Run(items, 3, false)
	assert.Len(t, processed, 3)
}

func Test_Run_Sorts(t *testing.T) {
	items := []*model.Tag{
		{Value: "cat", Score: 1},
		{Value: "dog", Score: 2},
		{Value: "foo", Score: 5},
		{Value: "bar", Score: 3},
		{Value: "bee", Score: 4},
	}
	processed := Run(items, 5, false)
	assert.Len(t, processed, 5)
	assert.Equal(t, "foo", processed[0].Value)
	assert.Equal(t, "bee", processed[1].Value)
	assert.Equal(t, "bar", processed[2].Value)
	assert.Equal(t, "dog", processed[3].Value)
	assert.Equal(t, "cat", processed[4].Value)
}

func Test_Run_DeDupes(t *testing.T) {
	items := []*model.Tag{
		{Value: "cat", Score: 5},
		{Value: "person", Score: 2},
		{Value: "people", Score: 5},
		{Value: "bar", Score: 3},
		{Value: "cats", Score: 1},
	}
	processed := Run(items, 5, false)
	assert.Len(t, processed, 3)
	assert.Equal(t, "people", processed[0].Value)
	assert.Equal(t, 7.0, processed[0].Score)
	assert.Equal(t, "cat", processed[1].Value)
	assert.Equal(t, 6.0, processed[1].Score)
	assert.Equal(t, "bar", processed[2].Value)
	assert.Equal(t, 3.0, processed[2].Score)
}

func Test_Run_IgnoresTFIDF_IfNoDocs(t *testing.T) {
	items := []*model.Tag{
		{Value: "cat", Score: 5},
	}
	processed := Run(items, 5, false)
	assert.Equal(t, 5.0, processed[0].Score)
}

func Test_Run_AppliesTFIDF_WithDocs(t *testing.T) {
	items := []*model.Tag{
		{Value: "cat", Score: 5, Docs: 1, DocsCount: 3},
	}
	processed := Run(items, 5, false)
	assert.Equal(t, 1.9684489712313906, processed[0].Score)
}

func Test_Run_AdjustsScores(t *testing.T) {
	items := []*model.Tag{
		{Value: "cat", Score: 5},
		{Value: "person", Score: 2},
		{Value: "people", Score: 5},
		{Value: "bar", Score: 3},
		{Value: "cats", Score: 1},
	}
	processed := Run(items, 5, true)
	assert.Len(t, processed, 3)
	assert.Equal(t, "people", processed[0].Value)
	assert.Equal(t, 1.0, processed[0].Score)
	assert.Equal(t, "cat", processed[1].Value)
	assert.Equal(t, 0.8571428571428571, processed[1].Score)
	assert.Equal(t, "bar", processed[2].Value)
	assert.Equal(t, 0.42857142857142855, processed[2].Score)
}
