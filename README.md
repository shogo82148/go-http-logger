[![test](https://github.com/shogo82148/go-http-logger/actions/workflows/test.yml/badge.svg)](https://github.com/shogo82148/go-http-logger/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/shogo82148/go-http-logger.svg)](https://pkg.go.dev/github.com/shogo82148/go-http-logger)

# go-http-logger

go-http-logger package is a logger for HTTP requests.
It is similar to "github.com/gorilla/handlers".LoggingHandler, but more flexible.

## SYNOPSIS

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shogo82148/go-http-logger"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func Logger(l httplogger.ResponseLog, r *http.Request) {
	log.Println("size:", l.RequestSize())
	log.Println("status:", l.Status())
	log.Println("method:", r.Method)
	log.Println("request uri:", r.RequestURI)
}

func main() {
	h := httplogger.LoggingHandler(httplogger.LoggerFunc(Logger), http.HandlerFunc(Handler))
	http.Handle("/", h)
	http.ListenAndServe(":8000", nil)
}
```

## LICENSE

This software is released under the MIT License, see LICENSE.md.
