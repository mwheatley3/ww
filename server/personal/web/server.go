package web

import (
	"fmt"
	"html/template"

	"github.com/Sirupsen/logrus"
	"github.com/mwheatley3/ww/server/auth"
	"github.com/mwheatley3/ww/server/httpserver"
	"github.com/mwheatley3/ww/server/personal/api/api"
)

const (
	cookieName = "wheatley"
)

var reactConfig = httpserver.ReactConfig{
	ClientRoutes: []string{
		"/",
		"/login",
		"/client/*splat",
	},
	Title:      "Working Wheatleys",
	AssetsPath: "js/public/personal",
}

// Config is the server config for the personal
// web server
type Config struct {
	httpserver.Config
	Addr   string
	Cookie struct {
		HashKey  []byte
		BlockKey []byte
	}
}

// NewServer returns a new personal server
func NewServer(l *logrus.Logger, s api.Service, c Config) *Server {

	rc := reactConfig
	rc.Config = c.Config
	rc.HeadExtra = template.HTML(fmt.Sprintf(`
<script type="text/javascript">
	window.ENV = {
		API_BASE_URL: "%s/api"
	};
</script>
	`, c.Addr))

	srv := &Server{
		server:  httpserver.NewReactServer(rc, l),
		logger:  l,
		service: s,
		cookie:  auth.NewCookieStore(cookieName, c.Cookie.HashKey, c.Cookie.BlockKey),
		addr:    c.Addr,
		conf:    c,
	}

	srv.server.Use(srv.authMiddleware())

	srv.server.Register("POST", "/api/auth", httpserver.HandlerFunc(srv.login))
	srv.server.Register("DELETE", "/api/auth", httpserver.HandlerFunc(srv.logout))
	srv.server.Register("GET", "/api/users/:userID", httpserver.HandlerFunc(srv.getUser))

	return srv
}

// Server is a react portal server
type Server struct {
	server  *httpserver.Server
	logger  *logrus.Logger
	service api.Service
	cookie  *auth.CookieStore
	addr    string
	conf    Config
}

// Listen starts the http server
func (s *Server) Listen() error {
	if err := s.service.Init(); err != nil {
		return err
	}

	return s.server.Listen()
}
