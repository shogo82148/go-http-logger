package httplogger

import (
	"testing"
)

func BenchmarkWrap(b *testing.B) {
	rw := &responseWriter{
		rw: &responseWriter{},
	}
	for i := 0; i < b.N; i++ {
		wrap(rw)
	}
}
