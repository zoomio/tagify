package tagify

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkTagify(b *testing.B) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port)))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Run(context.TODO(), Source(fmt.Sprintf("http://localhost:%d", port)),
			TargetType(HTML), Limit(5), NoStopWords(true))
		if err != nil {
			break
		}
	}
}
