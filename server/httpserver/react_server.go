package httpserver

import (
	"html/template"
	"net/http"

	"context"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

const (
	defaultReactAssetsPath   = "js/public"
	defaultReactAssetsPrefix = "/assets/"
	defaultJSIndex           = "main.js"
)

// ReactConfig paramaterizes the react http server
type ReactConfig struct {
	Config                     // httpserver.Config
	ClientRoutes []string      // a list of routes that should be handled in the client
	AssetsPath   string        // the fs path where the public assets dir is found
	AssetsPrefix string        // the url prefix to direct towards the fs path
	Title        string        // html title
	JSIndex      string        // js script name (assumed to be in the assets dir)
	HeadExtra    template.HTML // any extra html to shove into the head
	BodyExtra    template.HTML // extra html to shove into the body
}

// NewReactServer returns a new httpserver.Server that is set up
// for running a single page react app.
func NewReactServer(conf ReactConfig, l *logrus.Logger) *Server {
	if conf.AssetsPrefix == "" {
		conf.AssetsPrefix = defaultReactAssetsPrefix
	}

	if conf.AssetsPath == "" {
		conf.AssetsPath = defaultReactAssetsPath
	}

	if conf.JSIndex == "" {
		conf.JSIndex = defaultJSIndex
	}

	s := New(l, conf.Config)

	for _, r := range conf.ClientRoutes {
		s.Register("GET", r, indexHandler(conf, l))
	}

	s.Register("GET", conf.AssetsPrefix+"*splat", fsHandler(conf.AssetsPrefix, conf.AssetsPath))

	return s
}

func fsHandler(prefix, path string) Handler {
	fs := http.StripPrefix(prefix, http.FileServer(http.Dir(path)))

	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fs.ServeHTTP(w, r)
	})
}

func indexHandler(conf ReactConfig, l *logrus.Logger) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := indexTmpl.Execute(w, conf); err != nil {
			l.Errorf("React template error: %s", err)
		}
	})
}

var indexTmpl = template.Must(template.New("").Parse(`
<html>
	<head>
		<title>{{ .Title }}</title>
		{{ .HeadExtra }}
	</head>
	<body>
		{{ .BodyExtra }}
		<div id="react-doc"></div>
		<script src="{{ .AssetsPrefix }}{{ .JSIndex }}"></script>
	</body>
</html>
`))
