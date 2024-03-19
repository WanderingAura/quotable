package main

import (
	"flag"
	"net/http"
)

type application struct {
	config config
	logger string
}

type config struct {
	port string
}

func main() {
	var config config
	flag.StringVar(&config.port, "port", "4000", "server deploy port")

	app := &application{}

	srv := http.Server{
		Addr:    ":" + config.port,
		Handler: app.routes(),
	}

	srv.ListenAndServe()
}
