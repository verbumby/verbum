package chttp

import "net/http"

type middlewareFunc func(HandlerFunc) HandlerFunc

// ContentTypeJSONMiddleware http middleware to add conten-type response header
func ContentTypeJSONMiddleware(f HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, ctx *Context) error {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return f(w, ctx)
	}
}

func CacheControlImmutableMiddleware(f HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, ctx *Context) error {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		return f(w, ctx)
	}
}
