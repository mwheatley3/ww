package httpserver

import (
	"context"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Swapper is an http middleware that allows for runtime switching of
// the underlying handler
type Swapper struct {
	Handler Handler
	l       sync.RWMutex
}

// Swap sets the underlying handler to h and returns the old
// handler
func (s *Swapper) Swap(h Handler) Handler {
	s.l.Lock()
	old := s.Handler
	s.Handler = h
	s.l.Unlock()

	return old
}

// Handle satisfies the Handler interface
func (s *Swapper) Handle(c context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.l.RLock()
	h := s.Handler
	s.l.RUnlock()

	if h == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	h.Handle(c, w, r, p)
}
