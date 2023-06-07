package chttp

import "net/http"

// Context http request context
type Context struct {
	R *http.Request
	P *Principal
}
