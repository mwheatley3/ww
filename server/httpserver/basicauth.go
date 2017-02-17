package httpserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// BasicAuth is a middleware that provides basic auth validation
func BasicAuth(username, password string) MiddlewareFunc {
	return MiddlewareFunc(
		func(method, path string, h Handler) Handler {
			if username == "" {
				return h
			}

			return HandlerFunc(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				if user, pw, ok := r.BasicAuth(); ok && user == username && pw == password {
					h.Handle(w, r, ps)
					return
				}

				w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			})
		})
}
