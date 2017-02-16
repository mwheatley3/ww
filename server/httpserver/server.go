package httpserver

import (
	"net/http"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/facebookgo/httpdown"
	"github.com/julienschmidt/httprouter"
)

const (
	defaultStopTimeout = 5 * time.Second
	defaultKillTimeout = 10 * time.Second
)

// Config for the http server
type Config struct {
	ListenAddr   string
	StopTimeout  time.Duration
	KillTimeout  time.Duration
	OAuthEnabled bool
}

// A Route represents a Handler that has been
// registered with a Server
type Route struct {
	Method     string
	Path       string
	Handler    Handler
	Middleware []Middleware
}

// New returns a new server from a Config
func New(l *logrus.Logger, conf Config) *Server {
	s := &Server{
		Router: httprouter.New(),
		Config: conf,
		mws:    []Middleware{},
		routes: map[string]map[string]*Route{},
	}
	return s
}

// Server wraps a httpdown.Server and a httprouter.Router
// The Router is public so consumers of this package can
// access it directly.  Use of the server is wrapped
// through the Start, Stop and Wait methods
type Server struct {
	httpdown.Server
	Config  Config
	Router  *httprouter.Router
	mws     []Middleware
	routes  map[string]map[string]*Route
	startL  sync.Mutex
	routesL sync.Mutex
}

// Start starts the server and immediately
// returns.
func (s *Server) Start() error {
	s.startL.Lock()
	defer s.startL.Unlock()

	// we've already started
	if s.Server != nil {
		return nil
	}

	http := &http.Server{
		Addr:    s.Config.ListenAddr,
		Handler: s.Router,
	}

	conf := httpdown.HTTP{
		StopTimeout: s.Config.StopTimeout,
		KillTimeout: s.Config.KillTimeout,
	}

	if conf.StopTimeout == 0 {
		conf.StopTimeout = defaultStopTimeout
	}

	if conf.KillTimeout == 0 {
		conf.KillTimeout = defaultKillTimeout
	}

	server, err := conf.ListenAndServe(http)

	if err != nil {
		return err
	}

	s.Server = server
	return nil
}

// Listen is a shortcut for Start + Wait
// basically this is a blocking version of Start
// and provides familiar semantics (same as
// http.Server.ListenAndServe)
func (s *Server) Listen() error {
	defer s.closeMiddleware()
	if err := s.Start(); err != nil {
		return err
	}

	return s.Wait()
}

func (s *Server) closeMiddleware() {
	for _, mw := range s.mws {
		cmw, yep := mw.(Closeable)
		if yep {
			cmw.Close()
		}
	}
}

// Use adds the provided middleware to the server global middleware
func (s *Server) Use(mws ...Middleware) {
	s.mws = append(s.mws, mws...)
}

// Lookup returns a registered Route
func (s *Server) Lookup(method, path string) *Route {
	s.routesL.Lock()
	defer s.routesL.Unlock()

	routes, ok := s.routes[method]
	if !ok {
		return nil
	}

	return routes[path]
}

// Register adds a new route to the router with the optionally provided middleware
func (s *Server) Register(method, path string, h Handler, mws ...Middleware) *Route {
	s.routesL.Lock()
	defer s.routesL.Unlock()

	// global middleware
	for _, mw := range s.mws {
		h = mw.GetHandler(method, path, h)
	}

	// per route middleware
	for _, mw := range mws {
		h = mw.GetHandler(method, path, h)
	}

	if _, ok := s.routes[method]; !ok {
		s.routes[method] = map[string]*Route{}
	}

	s.routes[method][path] = &Route{
		Method:     method,
		Path:       path,
		Handler:    h,
		Middleware: mws,
	}

	s.Router.Handle(method, path, wrapHandler(h))
	return s.routes[method][path]
}
