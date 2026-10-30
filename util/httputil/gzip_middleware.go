package httputil

import (
	"net/http"
	"strings"
)

// GZIPMiddleware creates Middleware which encodes the response using gzip if the request accepts the gzip encoding.
// If the request does not accept gzip, this middleware has no effect. If the response is gzipped, headers will be set
// accordingly.
func GZIPMiddleware() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next(w, r)
				return
			}

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			w.Header().Del("Content-Length")

			gzipResponseWriter := NewGZIPResponseWriter(w)
			defer func() {
				_ = gzipResponseWriter.Close()
			}()

			next(gzipResponseWriter, r)
		}
	}
}
