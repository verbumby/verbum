package handlers

import "net/http"

// ToHTTPS redirects user to https
func ToHTTPS(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		target := "https://" + req.Host + req.URL.Path
		if len(req.URL.RawQuery) > 0 {
			target += "?" + req.URL.RawQuery
		}
		http.Redirect(w, req, target, http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, req)
	}
}
