package main

import (
	"net/http"

	"github.com/acoshift/kou/pkg/app"
)

func main() {
	http.ListenAndServe(":8080", app.Handler())
}
