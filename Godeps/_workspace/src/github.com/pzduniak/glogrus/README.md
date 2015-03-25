# goji/glogrus [![GoDoc](https://godoc.org/github.com/goji/glogrus?status.png)](https://godoc.org/github.com/goji/glogrus)

glogrus provides structured logging via logrus for Goji. 

## Example


```go

package main

import(
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web/middleware"
	"github.com/goji/glogrus"
	"github.com/Sirupsen/logrus"
)

func main() {
	goji.Abandon(middleware.Logger)

	logr := logrus.New()
	logr.Formatter = new(logrus.JSONFormatter)
	goji.Use(glogrus.NewGlogrus(logr, "my-app-name"))

	goji.Get("/ping", yourHandler)
	goji.Serve()
}

```

- - -


#### Looking for hierarchical structured logging?
[slog](https://github.com/zenazn/slog) and [lunk](https://github.com/codahale/lunk) looks interesting.