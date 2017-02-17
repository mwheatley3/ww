package httpserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// A Handler is an http.Handler with some additional functionality
type Handler interface {
	Handle(http.ResponseWriter, *http.Request, httprouter.Params)
}

// HandlerFunc is a func that satisfies Handler
type HandlerFunc func(http.ResponseWriter, *http.Request, httprouter.Params)

// Handle satisfies the Handler interface
func (fn HandlerFunc) Handle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fn(w, r, p)
}

func wrapHandler(h Handler) httprouter.Handle {
	return h.Handle
}

// A Middleware decorates a Handle with additional behavior.  Middleware is applied
// at route registration.
type Middleware interface {
	GetHandler(method, path string, h Handler) Handler
}

// MiddlewareFunc declares the signature for the wrapper func
type MiddlewareFunc func(method, path string, h Handler) Handler

// GetHandler Wraps anonymous GetHandlers for use as Middlewares
func (fn MiddlewareFunc) GetHandler(method, path string, h Handler) Handler {
	return fn(method, path, h)
}

// A Closeable may have other services which need to close when the Server closes
type Closeable interface {
	Close()
}
