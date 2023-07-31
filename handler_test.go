package httplogger

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHijack(t *testing.T) {
	var logged bool
	originalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusOK)
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error(err)
		}
		if !logged {
			t.Error("want logged, but not")
		}
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	})

	loggingHandler := LoggingHandler(LoggerFunc(func(l Attrs, r *http.Request) {
		if l.Status() != http.StatusSwitchingProtocols {
			t.Errorf("unexpected status code: %d", l.Status())
		}
		logged = true
	}), originalHandler)

	ts := httptest.NewServer(loggingHandler)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
