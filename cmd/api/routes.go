package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/quotes", app.listQuoteHandler)
	router.HandlerFunc(http.MethodGet, "/login", app.userLoginHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.userLoginPostHandler)
	router.HandlerFunc(http.MethodGet, "/quotes/:id", app.quoteHandler)
	router.HandlerFunc(http.MethodGet, "/users/:user_id/quotes/:quote_id", app.userQuoteHandler)

	return router
}
