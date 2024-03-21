package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/WanderingAura/quotable/internal/mailer"
	"github.com/WanderingAura/quotable/internal/models"
)

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   models.Models
	mailer   mailer.Mailer
	wg       sync.WaitGroup
}

type config struct {
	port    int
	env     string
	debug   bool
	logPath string
	db      struct {
		dsn                string
		maxOpenConnections int
		maxIdleConnections int
		maxIdleDuration    string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 4000, "API server port")
	flag.StringVar(&config.env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&config.debug, "debug", false, "debug mode")
	flag.StringVar(&config.logPath, "log-path", "/apps/quotable/api/logs", "File to write error logs in (absolute path)")
	flag.StringVar(&config.db.dsn, "dsn", "", "Postgres database source name")
	flag.IntVar(&config.db.maxOpenConnections, "db-max-open-conns", 25, "Postgres max open connections")
	flag.IntVar(&config.db.maxIdleConnections, "db-max-idle-conns", 25, "Postgres max idle connections")
	flag.StringVar(&config.db.maxIdleDuration, "db-max-idle-time", "15m", "Postgres max connection idle time")
	flag.Float64Var(&config.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&config.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&config.limiter.enabled, "limiter-enable", true, "Enable rate limiter")

	flag.Parse()

	logFile, err := os.Open(config.logPath)
	if err != nil {
		log.Fatal(err)
	}
	errorLog := log.New(logFile, "[ERROR]\t", log.Ldate|log.Ltime)
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(config)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		config:   config,
		errorLog: errorLog,
		infoLog:  infoLog,
		models:   models.New(db),
	}

	srv := http.Server{
		Addr:     fmt.Sprintf(":%d", config.port),
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	srv.ListenAndServe()
}

func openDB(config config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.db.maxOpenConnections)
	db.SetMaxIdleConns(config.db.maxIdleConnections)

	duration, err := time.ParseDuration(config.db.maxIdleDuration)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
