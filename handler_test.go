package httplogger_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/shogo82148/go-http-logger"
)

func ExampleLoggingHandler() {
	originalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	loggingHandler := httplogger.LoggingHandler(httplogger.LoggerFunc(func(l httplogger.ResponseLog, r *http.Request) {
		fmt.Println("size:", l.Size())
		fmt.Println("status:", l.Status())
		fmt.Println("method:", r.Method)
		fmt.Println("request uri:", r.RequestURI)
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
}
