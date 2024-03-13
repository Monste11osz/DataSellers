package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trc := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trc)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
