package httpserver

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// A Handler is an http.Handler with some additional functionality
type Handler interface {
	Handle(context.Context, http.ResponseWriter, *http.Request, httprouter.Params)
}

// HandlerFunc is a func that satisfies Handler
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request, httprouter.Params)

// Handle satisfies the Handler interface
func (fn HandlerFunc) Handle(c context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fn(c, w, r, p)
}

func wrapHandler(h Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		h.Handle(context.Background(), w, r, p)
	}
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
