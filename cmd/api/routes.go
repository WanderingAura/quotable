package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)
	router.HandlerFunc(http.MethodGet, "/discover", app.listQuoteHandler)
	router.HandlerFunc(http.MethodGet, "/login", app.userLoginHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.userLoginPostHandler)
	router.HandlerFunc(http.MethodGet, "/discover/:id", app.quoteHandler)

	return router

}
