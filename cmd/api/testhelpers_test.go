package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
)

type testServer struct {
	*httptest.Server
}

func mockServer(routes http.Handler) *testServer {
	server := httptest.NewServer(routes)
	return &testServer{server}
}

func mockApp() *application {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	return &application{
		logger: &logger,
	}
}

func (ts *testServer) get(t *testing.T, url string) (int, http.Header, string) {
	response, err := ts.Client().Get(ts.URL + url)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return response.StatusCode, response.Header, string(body)
}

func middlewareResponse(t *testing.T, middleware func(http.Handler) http.Handler, mockHandler http.Handler) (int, http.Header, string) {
	responseRecorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	middleware(mockHandler).ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return response.StatusCode, response.Header, string(body)
}