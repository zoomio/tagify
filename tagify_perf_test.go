package tagify

import (
	"fmt"
	"testing"
	
	"github.com/stretchr/testify/assert"

	"github.com/zoomio/tagify/processor"
)

func BenchmarkTagify(b *testing.B) {
	b.StopTimer()
	defer stopServer(startServer(fmt.Sprintf(":%d", port)))

	b.StartTimer()

	var tags []*processor.Tag
	var err error

	for i := 0; i < b.N; i++ {
		tags, err = GetTags(fmt.Sprintf("http://localhost:%d", port), HTML, 5, false, true)
		if err != nil {
			break
		}
	}

	b.StopTimer()

	assert.Nil(b, err)
	assert.Len(b, tags, 5)
	assert.Equal(b, []string{"him", "andread", "befell", "boy", "cakes"}, ToStrings(tags))
}