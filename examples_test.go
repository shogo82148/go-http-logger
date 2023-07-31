package httplogger_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	httplogger "github.com/shogo82148/go-http-logger"
)

func ExampleLoggingHandler() {
	originalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "Hello World")
	})

	loggingHandler := httplogger.LoggingHandler(httplogger.LoggerFunc(func(l httplogger.Attrs, r *http.Request) {
		fmt.Println("size:", l.ResponseSize())
		fmt.Println("status:", l.Status())
		fmt.Println("method:", r.Method)
		fmt.Println("request uri:", r.RequestURI)
		fmt.Println("content type:", l.Header().Get("Content-Type"))
	}), originalHandler)

	ts := httptest.NewServer(loggingHandler)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Output:
	// size: 11
	// status: 200
	// method: GET
	// request uri: /
	// content type: text/plain
}
