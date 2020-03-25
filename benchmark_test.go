package httplogger

import (
	"net/http/httptest"
	"testing"
)

func BenchmarkWrap(b *testing.B) {
	rw := &responseWriter{
		rw: httptest.NewRecorder(),
	}
	for i := 0; i < b.N; i++ {
		wrap(rw)
	}
}
