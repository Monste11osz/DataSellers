package main

import (
	"net/http"
)

func (app *application) router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/signIn", app.signIn)
	mux.HandleFunc("/us/authentication", app.authentication)
	mux.HandleFunc("/us/form/inputFile", RequireCookie(app.start))
	mux.HandleFunc("/us/form/products/search", RequireCookie(app.search))
	return mux
}
