package app

import (
	"net/http"

	"github.com/acoshift/header"
	"github.com/acoshift/kou/pkg/util"
	"github.com/acoshift/middleware"
)

// Handler gets kou's handler
func Handler() http.Handler {
	m := http.NewServeMux()
	m.Handle("/", mux)
	m.Handle("/kou-admin/", configMux)
	m.Handle("/kou-content/", http.StripPrefix("/kou-content", util.FileServer("kou-content")))
	return middleware.Chain(middlewares...)(m)
}

// NotFound serves not found template
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(header.XContentTypeOptions, "nosniff")
	ExecuteTemplateWithCode(w, r, http.StatusNotFound, "not-found", nil)
}

// InternalServerError serves internal server error template
func InternalServerError(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set(header.XContentTypeOptions, "nosniff")
	ExecuteTemplateWithCode(w, r, http.StatusInternalServerError, "internal-server-error", data)
}
