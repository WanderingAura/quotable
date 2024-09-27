package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/WanderingAura/quotable/internal/assert"
)

func TestVersionCheckHandler(t *testing.T) {

	testEnvs := []string{
		"release",
		"debug",
		"development",
	}

	app := mockApp()

	ts := mockServer(app.routes())
	defer ts.Close()

	for _, env := range testEnvs {
		app.config.env = env

		statusCode, _, body := ts.get(t, "/v1/version")
		assert.Equal(t, statusCode, http.StatusOK)

		expectedProperty := fmt.Sprintf("\"environment\": \"%s\"", env)
		assert.StringContains(t, body, expectedProperty)
	}
}
