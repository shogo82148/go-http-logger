//go:build go1.21
// +build go1.21

package httplogger_test

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"

	httplogger "github.com/shogo82148/go-http-logger"
)

func ExampleLoggingHandler_slog() {
	originalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "Hello World")
	})

	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove time.
		if a.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}
		// Remove unstable attributes for testing.
		if a.Key == "received_time" || a.Key == "request_time" || a.Key == "vhost" || a.Key == "host" {
			return slog.Attr{}
		}
		return a
	}

	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replace}))
	logger := httplogger.NewSlogLogger(slog.LevelInfo, "test", slogger)
	loggingHandler := httplogger.LoggingHandler(logger, originalHandler)

	ts := httptest.NewServer(loggingHandler)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Output:
	// level=INFO msg=test forwardedfor="" user="" req="GET / HTTP/1.1" method=GET uri=/ protocol=HTTP/1.1 status=200 sent_bytes=11 received_bytes=0 referer="" ua=Go-http-client/1.1
}
