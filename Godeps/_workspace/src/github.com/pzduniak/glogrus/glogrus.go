package glogrus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

// glogrus is a middleware handler that logs the
// request and response in a structured way
type glogrus struct {
	h    http.Handler
	c    *web.C
	l    *logrus.Logger
	name string
}

func (glogr glogrus) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	reqID := middleware.GetReqID(*glogr.c)
	glogr.l.WithFields(logrus.Fields{
		"req_id": reqID,
		"uri":    req.RequestURI,
		"method": req.Method,
		"remote": req.RemoteAddr,
	}).Info("New request incoming")
	lresp := wrapWriter(resp)

	glogr.h.ServeHTTP(lresp, req)
	lresp.maybeWriteHeader()

	latency := float64(time.Since(start)) / float64(time.Millisecond)

	glogr.l.WithFields(logrus.Fields{
		"req_id":  reqID,
		"status":  lresp.status(),
		"method":  req.Method,
		"uri":     req.RequestURI,
		"remote":  req.RemoteAddr,
		"latency": fmt.Sprintf("%6.4f ms", latency),
		"app":     glogr.name,
	}).Info("Finished serving the request")
}

// NewGlogrus allows you to configure a goji middleware that logs all requests and responses
// using the structured logger logrus. It takes the logrus instance and the name of the app
// as the parameters and returns a middleware of type "func(c *web.C, http.Handler) http.Handler"
//
// Example:
//
//		package main
//
//		import(
//			"github.com/zenazn/goji"
//			"github.com/zenazn/goji/web/middleware"
//			"github.com/goji/glogrus"
//			"github.com/Sirupsen/logrus"
//		)
//
//		func main() {
//			goji.Abandon(middleware.Logger)
//
//			logr := logrus.New()
//			logr.Formatter = new(logrus.JSONFormatter)
//			goji.Use(glogrus.NewGlogrus(logr, "my-app-name"))
//
//			goji.Get("/ping", yourHandler)
//			goji.Serve()
//		}
//
func NewGlogrus(l *logrus.Logger, name string) func(*web.C, http.Handler) http.Handler {
	fn := func(c *web.C, h http.Handler) http.Handler {
		return glogrus{h: h, c: c, l: l, name: name}
	}
	return fn
}
