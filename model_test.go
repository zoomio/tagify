package tagify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zoomio/tagify/processor/model"
)

const txt = "Some random text to test Tagify model"

func Test_ForEach(t *testing.T) {
	res, err := Run(
		context.TODO(),
		Content(txt),
		TargetType(Text),
		Limit(3),
		NoStopWords(true),
	)
	assert.Nil(t, err)

	var count int

	it := func(i int, tag *model.Tag) {
		count++
	}

	res.ForEach(it)

	assert.Equal(t, 3, count)
}
