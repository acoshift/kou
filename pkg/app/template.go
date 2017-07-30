package app

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/acoshift/header"
)

// ParseTemplate parses template from text
func ParseTemplate(text string) (*template.Template, error) {
	return template.New("").
		Funcs(tfuncs).
		Parse(text)
}

// ExecuteTemplateWithCode executes a template from given name with http code
func ExecuteTemplateWithCode(w http.ResponseWriter, r *http.Request, code int, name string, data interface{}) {
	t, ok := templates[name]
	if !ok {
		logger.Printf("template/exec: template \"%s\" not found\n", name)
		http.NotFound(w, r)
	}

	w.Header().Set(header.ContentType, "text/html; charset=utf-8")
	w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate, max-age=0")
	w.WriteHeader(code)
	buf := bytes.Buffer{}
	err := t.Execute(&buf, data)
	if err != nil {
		logger.Printf("template/exec: execute error; %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = minifier.Minify("text/html", w, &buf)
	if err != nil {
		logger.Printf("template/exec: minify error; %v\n", err)
		return
	}
}

// ExecuteTemplate executes a template from given name
func ExecuteTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	ExecuteTemplateWithCode(w, r, http.StatusOK, name, data)
}

// ServeTemplate returns handler to execute template
func ServeTemplate(name string, data interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ExecuteTemplate(w, r, name, data)
	})
}
