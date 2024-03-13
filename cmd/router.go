package main

import (
	"net/http"
)

func (app *application) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/authentication", app.authentication)
	mux.HandleFunc("/us/inputFile", RequireCookie(app.start))
	mux.HandleFunc("/us/products/search", RequireCookie(app.search))
	return mux
}
