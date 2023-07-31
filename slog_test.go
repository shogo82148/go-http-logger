//go:build go1.21
// +build go1.21

package httplogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSlog(t *testing.T) {
	buf := &bytes.Buffer{}
	jsonHandler := slog.NewJSONHandler(buf, nil)
	jsonLogger := slog.New(jsonHandler)
	logger := NewSlogLogger(slog.LevelInfo, "test", jsonLogger)
	originalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello World")
	})
	loggingHandler := LoggingHandler(logger, originalHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-Agent", "test-agent/1.0")
	req.SetBasicAuth("test-user", "test-password")
	rw := httptest.NewRecorder()
	loggingHandler.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("unexpected status code: %d", rw.Code)
	}

	var log map[string]any
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Error(err)
	}

	if _, ok := log["time"].(string); !ok {
		t.Errorf("time is not logged")
	}
	if log["level"] != "INFO" {
		t.Errorf("unexpected level: got %q, want %q", log["level"], "INFO")
	}
	if _, ok := log["received_time"].(string); !ok {
		t.Errorf("received_time is not logged")
	}
	if log["host"] != "192.0.2.1:1234" {
		t.Errorf("unexpected host: got %q, want %q", log["host"], "192.0.2.1:1234")
	}
	if log["forwardedfor"] != "" {
		t.Errorf("unexpected xforwardedfor: got %q, want %q", log["xforwardedfor"], "")
	}
	if log["user"] != "test-user" {
		t.Errorf("unexpected user: got %q, want %q", log["user"], "test-user")
	}
	if log["req"] != "GET / HTTP/1.1" {
		t.Errorf("unexpected req: got %q, want %q", log["req"], "GET / HTTP/1.1")
	}
	if log["method"] != "GET" {
		t.Errorf("unexpected method: got %q, want %q", log["method"], "GET")
	}
	if log["uri"] != "/" {
		t.Errorf("unexpected uri: got %q, want %q", log["uri"], "/")
	}
	if log["protocol"] != "HTTP/1.1" {
		t.Errorf("unexpected protocol: got %q, want %q", log["protocol"], "HTTP/1.1")
	}
	if log["status"] != float64(200) {
		t.Errorf("unexpected status: got %f, want %f", log["status"], float64(200))
	}
	if log["sent_bytes"] != float64(11) {
		t.Errorf("unexpected size: got %f, want %f", log["size"], float64(11))
	}
	if log["received_bytes"] != float64(0) {
		t.Errorf("unexpected received_bytes: got %f, want %f", log["received_bytes"], float64(0))
	}
	if log["referer"] != "" {
		t.Errorf("unexpected referer: got %q, want %q", log["referer"], "")
	}
	if log["ua"] != "test-agent/1.0" {
		t.Errorf("unexpected ua: got %q, want %q", log["ua"], "test-agent/1.0")
	}
	if log["host"] != "192.0.2.1:1234" {
		t.Errorf("unexpected host: got %q, want %q", log["host"], "")
	}
	if _, ok := log["request_time"].(float64); !ok {
		t.Errorf("request_time is not logged")
	}
}
