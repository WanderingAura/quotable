package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	// TODO: create an errors page and use it here
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.ErrorTitle = "404 page not found"
	data.ErrorContent = "Please try searching for the page below"
	// TODO: render the error.html template with the data
}

func (app *application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("Method not allowed")
	app.serverErrorResponse(w, r, err)
}
