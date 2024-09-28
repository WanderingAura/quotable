package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

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

func TestRateLimit(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	app := mockApp()
	app.config.limiter.enabled = true

	type LimiterTest struct {
		name  string
		rps   float64
		burst int
	}

	limiterTests := []LimiterTest{
		{
			name:  "same rps and burst",
			rps:   10.0,
			burst: 10,
		},
		{
			name:  "rps < burst",
			rps:   10,
			burst: 40,
		},
		{
			name:  "rps > burst",
			rps:   13.4,
			burst: 5,
		},
	}

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
	})

	for _, test := range limiterTests {
		t.Run(test.name, func(t *testing.T) {
			app.config.limiter.burst = test.burst
			app.config.limiter.rps = test.rps
			next := app.rateLimit(mockHandler)
			for i := 0; i < test.burst-1; i++ {
				handlerResponse(t, next)
			}

			statusCode, _, body := handlerResponse(t, next)
			assert.Equal(t, statusCode, http.StatusOK)
			assert.StringContains(t, body, "OK")

			statusCode, _, body = handlerResponse(t, next)
			assert.Equal(t, statusCode, http.StatusTooManyRequests)
			assert.StringContains(t, body, "rate limit exceeded")

			timeToRefillAlmostMS := int(500 / test.rps)
			leftOverTime := int(500 / test.rps)
			const graceMS = 10

			// wait a small time
			time.Sleep(time.Duration(timeToRefillAlmostMS) * time.Millisecond)

			statusCode, _, body = handlerResponse(t, next)
			assert.Equal(t, statusCode, http.StatusTooManyRequests)
			assert.StringContains(t, body, "rate limit exceeded")

			time.Sleep(time.Duration(leftOverTime+graceMS) * time.Millisecond)

			statusCode, _, body = handlerResponse(t, next)
			assert.Equal(t, statusCode, http.StatusOK)
			assert.StringContains(t, body, "OK")
		})

	}

	limiterErrorTests := []LimiterTest{
		{
			name:  "negative rps",
			rps:   -3.1,
			burst: 5,
		},
		{
			name:  "negative burst",
			rps:   13.4,
			burst: -5,
		},
	}

	for _, test := range limiterErrorTests {
		t.Run(test.name, func(t *testing.T) {
			app.config.limiter.burst = test.burst
			app.config.limiter.rps = test.rps
			next := app.rateLimit(mockHandler)
			for i := 0; i < 10; i++ {
				statusCode, _, body := handlerResponse(t, next)
				assert.Equal(t, statusCode, http.StatusTooManyRequests)
				assert.StringContains(t, body, "rate limit exceeded")
			}
		})

	}
}
