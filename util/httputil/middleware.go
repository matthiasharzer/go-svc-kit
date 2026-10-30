package httputil

import (
	"net/http"
	"slices"
)

// A Middleware is a wrapper functions applied before the actual http.HandlerFunc used for manipulating behavior and
// altering requests / responses.
type Middleware = func(next http.HandlerFunc) http.HandlerFunc

// UseMiddleware wraps the handler with all middlewares. The order of middlewares will be respected. The first
// Middleware will be applied first to incoming requests. The handler last.
func UseMiddleware(middlewares []Middleware, handler http.HandlerFunc) http.HandlerFunc {
	reversedMiddlewares := make([]Middleware, len(middlewares))
	copy(reversedMiddlewares, middlewares)

	// the first middleware in the array should be applied first to the request, so it has to be the last one that
	// wraps the handler. Therefore, the middlewares are applied in reverse order
	slices.Reverse(reversedMiddlewares)

	finalHandler := handler
	for _, middleware := range reversedMiddlewares {
		finalHandler = middleware(finalHandler)
	}

	return finalHandler
}
