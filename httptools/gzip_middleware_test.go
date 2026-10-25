package httptools_test

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matthiasharzer/service-kit/httptools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGZIPMiddleware(t *testing.T) {
	t.Run("serves gzip response if gzip is accepted", func(t *testing.T) {
		handler := httptools.UseMiddleware(
			[]httptools.Middleware{
				httptools.GZIPMiddleware(),
			},
			func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("hello world"))
				assert.NoError(t, err)
			},
		)

		request, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)
		request.Header.Set("Accept-Encoding", "gzip")

		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, request)
		require.Equal(t, http.StatusOK, recorder.Code)

		assert.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))

		gzipReader, err := gzip.NewReader(recorder.Body)
		require.NoError(t, err)
		gzipBody, err := io.ReadAll(gzipReader)
		require.NoError(t, err)
		assert.Equal(t, "hello world", string(gzipBody))
	})

	t.Run("does not serve gzip if it is not accepted", func(t *testing.T) {
		handler := httptools.UseMiddleware(
			[]httptools.Middleware{
				httptools.GZIPMiddleware(),
			},
			func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("hello world"))
				assert.NoError(t, err)
			},
		)

		request, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, request)
		require.Equal(t, http.StatusOK, recorder.Code)
		assert.NotEqual(t, "gzip", recorder.Header().Get("Content-Encoding"))

		body, err := io.ReadAll(recorder.Result().Body)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", string(body))
	})
}
