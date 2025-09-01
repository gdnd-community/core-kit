package log

import (
	"testing"

	"github.com/gdnd-community/core-kit/pkg/meta"
)

func Test_Logger(t *testing.T) {

	metas := meta.Discover("CoreKit Test", "1.0.0", "Local")
	Init("info", WithMetadata(metas))

	Info("Test with microservice core-kit", map[string]interface{}{
		"module": "test",
	})
}

func Benchmark_Logger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("test", map[string]any{"as": "as"})
	}
}
