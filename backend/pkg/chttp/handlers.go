package chttp

import (
	"log"
	"net/http"
)

// HandlerFunc app http handler func type
type HandlerFunc func(http.ResponseWriter, *Context) error

// MakeHandler creates http.HandlerFunc of chttp.handlerFunc and specified middlewares
func MakeHandler(f HandlerFunc, middlewares ...middlewareFunc) http.HandlerFunc {
	f = elasticAccessLogMiddleware(f)
	for _, m := range middlewares {
		f = m(f)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{}
		ctx.R = r
		if err := f(w, ctx); err != nil {
			log.Printf("%s %s: error: %v", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// MakeAPIHandler creates api handler with content type and auth middlewares
func MakeAPIHandler(f HandlerFunc) http.HandlerFunc {
	return MakeHandler(f, ContentTypeJSONMiddleware, AuthMiddleware)
}
