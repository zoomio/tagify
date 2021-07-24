package tagify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zoomio/tagify/config"
	"github.com/zoomio/tagify/processor/model"
)

const txt = "Some random text to test Tagify model"

var ctx = context.TODO()

func Test_ForEach(t *testing.T) {
	res, err := Run(ctx,
		config.Content(txt),
		config.TargetType(config.Text),
		config.Limit(3),
		config.NoStopWords(true),
	)
	assert.Nil(t, err)

	var count int

	it := func(i int, tag *model.Tag) {
		count++
	}

	res.ForEach(it)

	assert.Equal(t, 3, count)
}
