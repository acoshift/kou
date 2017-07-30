package app

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/acoshift/middleware"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

var (
	mux         = http.NewServeMux()
	configMux   = http.NewServeMux()
	tfuncs      = make(template.FuncMap)
	logger      = log.New(os.Stderr, "", 0)
	middlewares = make([]middleware.Middleware, 0)
	templates   = make(map[string]*template.Template)
	minifier    = minify.New()
)

func init() {
	// init
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/javascript", js.Minify)

	// add default index template
	RegisterTemplate("index", template.Must(ParseTemplate(`Welcome to Kou`)))
	RegisterTemplate("not-found", template.Must(ParseTemplate(`404 Page Not Found Kou~!`)))
	RegisterTemplate("internal-server-error", template.Must(ParseTemplate(`500 Internal Server Error {{.}}`)))

	RegisterHandler("/", OnlyRootPath()(ServeTemplate("index", nil)))
}

// RegisterHandler registers handler into kou http mux
func RegisterHandler(pattern string, handler http.Handler) {
	logger.Printf("register/handler: add %s\n", pattern)
	mux.Handle(pattern, handler)
}

// RegisterFunc registers func into kou template func map
func RegisterFunc(name string, f interface{}) {
	if _, ok := tfuncs[name]; ok {
		logger.Printf("register/func: override %s\n", name)
	} else {
		logger.Printf("register/func: add %s\n", name)
	}
	tfuncs[name] = f
}

// RegisterConfig registers config handler into kou config handler
func RegisterConfig(pattern string, handler http.Handler) {
	logger.Printf("register/config: add %s\n", pattern)
	configMux.Handle(pattern, handler)
}

// RegisterTemplate registers template into kou templates
func RegisterTemplate(name string, t *template.Template) {
	if _, ok := templates[name]; ok {
		logger.Printf("register/templates: override %s\n", name)
	} else {
		logger.Printf("register/templates: add %s\n", name)
	}
	templates[name] = t
}

// RegisterMiddleware registers middleware to kou middleware
func RegisterMiddleware(name string, m middleware.Middleware) {
	logger.Printf("register/middleware: add %s\n", name)
	middlewares = append(middlewares, m)
}
