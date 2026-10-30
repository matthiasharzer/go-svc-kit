package httputil_test

import (
	"io/fs"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/matthiasharzer/service-kit/util/httputil"
	"github.com/stretchr/testify/assert"
)

func getTestFS() fs.FS {
	return fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte("index"),
		},
		"foo.html": &fstest.MapFile{
			Data: []byte("foo"),
		},
		"dist/bar.html": &fstest.MapFile{
			Data: []byte("bar"),
		},
		"bar/index.html": &fstest.MapFile{
			Data: []byte("bar index"),
		},
	}
}

func TestHandleStaticSite(t *testing.T) {
	t.Run("serves files in the filesystem", func(t *testing.T) {
		testFS := getTestFS()

		handler := httputil.HandleStaticSite(testFS)

		request := httptest.NewRequest("GET", "/index.html", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "index", recorder.Body.String())

		request = httptest.NewRequest("GET", "/foo.html", nil)
		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "foo", recorder.Body.String())

		request = httptest.NewRequest("GET", "/dist/bar.html", nil)
		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "bar", recorder.Body.String())

		request = httptest.NewRequest("GET", "/bar/index.html", nil)
		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "bar index", recorder.Body.String())
	})

	t.Run("serves index.html if path ends with /", func(t *testing.T) {
		testFS := getTestFS()

		handler := httputil.HandleStaticSite(testFS)

		request := httptest.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "index", recorder.Body.String())

		request = httptest.NewRequest("GET", "/bar/", nil)
		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "bar index", recorder.Body.String())
	})

	t.Run("serves index.html if file does not exist", func(t *testing.T) {
		testFS := getTestFS()

		handler := httputil.HandleStaticSite(testFS)

		request := httptest.NewRequest("GET", "/nonexistent.html", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "index", recorder.Body.String())
	})

	t.Run("returns 500 if index.html does not exist", func(t *testing.T) {
		testFS := fstest.MapFS{}

		handler := httputil.HandleStaticSite(testFS)

		request := httptest.NewRequest("GET", "/nonexistent.html", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, 500, recorder.Code)
	})
}
