package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/WanderingAura/quotable/internal/assert"
)

const TestServerErrorString = "the server encountered a problem and could not process your request"

func TestRecoverPanic(t *testing.T) {
	app := mockApp()

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("THIS HANDLER HAS PANICKED")
	})

	statusCode, header, body := middlewareResponse(t, app.recoverPanic, mockHandler)

	assert.Equal(t, statusCode, http.StatusInternalServerError)
	assert.Equal(t, header.Get("Connection"), "close")
	expectedContains := fmt.Sprintf("\"error\": \"%s\"", TestServerErrorString)
	assert.StringContains(t, body, expectedContains)
}
