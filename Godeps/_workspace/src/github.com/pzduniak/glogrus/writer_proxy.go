package glogrus

import (
	"bufio"
	"net"
	"net/http"
)

// wrapWriter returns a proxy that wraps ResponseWriter
func wrapWriter(w http.ResponseWriter) writerProxy {
	bw := basicWriter{ResponseWriter: w}
	return &bw
}

// writerProxy is a proxy that wraps ResponseWriter
type writerProxy interface {
	http.ResponseWriter
	maybeWriteHeader()
	status() int
}

// basicWriter holds the status code and a
// flag in addition to http.ResponseWriter
type basicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
}

// WriteHeader stores the status code and writes header
func (b *basicWriter) WriteHeader(code int) {
	if !b.wroteHeader {
		b.code = code
		b.wroteHeader = true
		b.ResponseWriter.WriteHeader(code)
	}
}

// Write writes the bytes and calls MaybeWriteHeader
func (b *basicWriter) Write(buf []byte) (int, error) {
	b.maybeWriteHeader()
	return b.ResponseWriter.Write(buf)
}

// maybeWriteHeader writes the header if it is not alredy set
func (b *basicWriter) maybeWriteHeader() {
	if !b.wroteHeader {
		b.WriteHeader(http.StatusOK)
	}
}

// status returns the status
func (b *basicWriter) status() int {
	return b.code
}

// unwrap returns the original http.ResponseWriter
func (b *basicWriter) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}

// Flush is used for the implementation of http.Flusher
func (b *basicWriter) Flush() {
	b.ResponseWriter.(http.Flusher).Flush()
}

// Hijack is used for the implementation of http.Hijacker
func (b *basicWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return b.ResponseWriter.(http.Hijacker).Hijack()
}
