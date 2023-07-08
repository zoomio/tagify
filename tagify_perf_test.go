package tagify

import (
	"context"
	"fmt"
	"os"
	"testing"
)

var (
	vergeHTML, _ = os.ReadFile("../../_resources_test/html/theverge.html")
)

func BenchmarkTagify(b *testing.B) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port), string(vergeHTML)))

	b.ResetTimer()

	ctx := context.TODO()

	for i := 0; i < b.N; i++ {
		_, err := Run(ctx,
			Source(fmt.Sprintf("http://localhost:%d", port)),
			TargetType(HTML),
			Limit(40),
			NoStopWords(true),
		)
		if err != nil {
			b.Fatal(err)
			break
		}
	}
}
