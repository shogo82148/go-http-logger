package httplogger_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	httplogger "github.com/shogo82148/go-http-logger"
)

func TestWriteHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	defer log.SetOutput(os.Stderr)

	originalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusNotFound) // want warning
		fmt.Fprint(w, "Hello World")
	})

	loggingHandler := httplogger.LoggingHandler(httplogger.LoggerFunc(func(l httplogger.Attrs, r *http.Request) {
		if l.Status() != http.StatusOK {
			t.Errorf("unexpected status code: %d", l.Status())
		}
	}), originalHandler)

	rw := httptest.NewRecorder()
	loggingHandler.ServeHTTP(rw, httptest.NewRequest("GET", "/", nil))
	if rw.Code != http.StatusOK {
		t.Errorf("unexpected status code: %d", rw.Code)
	}

	if !strings.Contains(buf.String(), "logger_test.go") {
		t.Errorf("unexpected log: %s", buf.String())
	}
}
