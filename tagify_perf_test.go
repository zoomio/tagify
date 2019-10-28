package tagify

import (
	"fmt"
	"testing"
)

func BenchmarkTagify(b *testing.B) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port)))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GetTags(fmt.Sprintf("http://localhost:%d", port), HTML, 5, false, true)
		if err != nil {
			break
		}
	}
}
