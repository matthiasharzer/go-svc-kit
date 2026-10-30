package httputil

import (
	"compress/gzip"
	"net/http"
)

// GZIPResponseWriter wraps an http.ResponseWriter with a gzip.Writer enabling gzip encoded http responses
type GZIPResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (g *GZIPResponseWriter) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

func (g *GZIPResponseWriter) Close() error {
	return g.writer.Close()
}

func NewGZIPResponseWriter(w http.ResponseWriter) *GZIPResponseWriter {
	gzipWriter := gzip.NewWriter(w)
	return &GZIPResponseWriter{
		ResponseWriter: w,
		writer:         gzipWriter,
	}
}
