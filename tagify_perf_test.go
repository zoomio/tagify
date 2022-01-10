package tagify

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkTagify(b *testing.B) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port), indexHTML))

	b.ResetTimer()

	ctx := context.TODO()

	for i := 0; i < b.N; i++ {
		_, err := Run(ctx,
			Source(fmt.Sprintf("http://localhost:%d", port)),
			TargetType(HTML),
			Limit(5),
			NoStopWords(true),
		)
		if err != nil {
			break
		}
	}
}
