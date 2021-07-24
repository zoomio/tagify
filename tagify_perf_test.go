package tagify

import (
	"context"
	"fmt"
	"testing"

	"github.com/zoomio/tagify/config"
)

func BenchmarkTagify(b *testing.B) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port)))

	b.ResetTimer()

	ctx := context.TODO()

	for i := 0; i < b.N; i++ {
		_, err := Run(ctx,
			config.Source(fmt.Sprintf("http://localhost:%d", port)),
			config.TargetType(config.HTML),
			config.Limit(5),
			config.NoStopWords(true),
		)
		if err != nil {
			break
		}
	}
}
