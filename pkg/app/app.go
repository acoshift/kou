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

type muxItem struct {
	pattern string
	handler http.Handler
}

var (
	mux         = make(map[string]http.Handler)
	configMux   = http.NewServeMux()
	tfuncs      = make(template.FuncMap)
	logger      = log.New(os.Stderr, "", 0)
	middlewares = make([]middleware.Middleware, 0)
	templates   = make(map[string]*template.Template)
	minifier    = minify.New()
)

const (
	kouModule = "kou"
)

func init() {
	// init
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/javascript", js.Minify)

	// add default index template
	RegisterTemplate(kouModule, "index", template.Must(ParseTemplate(`Welcome to Kou`)))
	RegisterTemplate(kouModule, "not-found", template.Must(ParseTemplate(`404 Page Not Found Kou~!`)))
	RegisterTemplate(kouModule, "internal-server-error", template.Must(ParseTemplate(`500 Internal Server Error {{.}}`)))

	RegisterHandler(kouModule, "/", OnlyRootPath()(ServeTemplate("index", nil)))
}

// RegisterHandler registers handler into kou http mux
func RegisterHandler(module, pattern string, handler http.Handler) {
	if _, ok := mux[pattern]; ok {
		logger.Printf("register/handler: [%s] override %s\n", module, pattern)
	} else {
		logger.Printf("register/handler: [%s] add %s\n", module, pattern)
	}
	mux[pattern] = handler
}

// RegisterFunc registers func into kou template func map
func RegisterFunc(module, name string, f interface{}) {
	if _, ok := tfuncs[name]; ok {
		logger.Printf("register/func: [%s] override %s\n", module, name)
	} else {
		logger.Printf("register/func: [%s] add %s\n", module, name)
	}
	tfuncs[name] = f
}

// RegisterConfig registers config handler into kou config handler
func RegisterConfig(module, pattern string, handler http.Handler) {
	logger.Printf("register/config: [%s] add %s\n", module, pattern)
	configMux.Handle(pattern, handler)
}

// RegisterTemplate registers template into kou templates
func RegisterTemplate(module, name string, t *template.Template) {
	if _, ok := templates[name]; ok {
		logger.Printf("register/templates: [%s] override %s\n", module, name)
	} else {
		logger.Printf("register/templates: [%s] add %s\n", module, name)
	}
	templates[name] = t
}

// RegisterMiddleware registers middleware to kou middleware
func RegisterMiddleware(module, name string, m middleware.Middleware) {
	logger.Printf("register/middleware: [%s] add %s\n", module, name)
	middlewares = append(middlewares, m)
}
