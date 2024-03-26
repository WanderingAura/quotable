package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/quotes", app.listQuotesHandler)
	router.HandlerFunc(http.MethodGet, "/tokens/auth", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/user/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodGet, "/quotes/:id", app.getQuoteHandler)
	router.HandlerFunc(http.MethodGet, "/users/:user_id/quotes", app.listUserQuotesHandler)

	return router
}
