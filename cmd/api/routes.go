package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/version", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/quotes", app.listQuotesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/auth", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/user/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/quotes/:quote_id", app.getQuoteHandler)
	router.HandlerFunc(http.MethodPost, "/v1/quotes", app.requireAuthenticatedUser(app.createQuoteHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/:user_id/quotes", app.requireAuthenticatedUser(app.listUserQuotesHandler))

	return app.rateLimit(app.authenticate(router))
}
