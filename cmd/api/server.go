package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	// TODO: create custom logger which logs errors and info into a file as well as logs info to stdout.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.errorLog.Writer(), "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1) // use a buffer size of 1 to avoid missing signals when quit is not ready to receive
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.infoLog.Printf("signal: %s\n shutting down server", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.infoLog.Printf("completing background tasks, port: %s", srv.Addr)

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.infoLog.Printf("starting server, port: %s, env: %s", srv.Addr, app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.infoLog.Printf("stopped server, port: %s", srv.Addr)

	return nil
}
