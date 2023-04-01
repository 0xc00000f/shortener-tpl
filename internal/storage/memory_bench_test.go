package storage

import (
	"testing"

	"go.uber.org/zap"
)

func BenchmarkMemoryStorage_Get(b *testing.B) {
	ms := NewMemoryStorage(zap.L())
	for i := 0; i < b.N; i++ {
		ms.Get(nil, "short")
	}
}
