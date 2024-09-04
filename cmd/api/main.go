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

	"github.com/WanderingAura/quotable/internal/data"
	"github.com/WanderingAura/quotable/internal/mailer"
)

var (
	buildTime string
	version   string
)

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Models
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
	flag.StringVar(&config.logPath, "log-path", "./logs/quotable.log", "File to write error logs in")
	flag.StringVar(&config.db.dsn, "db-dsn", "", "Postgres database source name")
	flag.IntVar(&config.db.maxOpenConnections, "db-max-open-conns", 25, "Postgres max open connections")
	flag.IntVar(&config.db.maxIdleConnections, "db-max-idle-conns", 25, "Postgres max idle connections")
	flag.StringVar(&config.db.maxIdleDuration, "db-max-idle-time", "15m", "Postgres max connection idle time")
	flag.Float64Var(&config.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&config.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&config.limiter.enabled, "limiter-enable", true, "Enable rate limiter")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	if config.db.dsn == "" {
		config.db.dsn = os.Getenv("QUOTABLE_DSN")
	}
	// logFile, err := os.Open(config.logPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	errorLog := log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime)
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(config)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		config:   config,
		errorLog: errorLog,
		infoLog:  infoLog,
		models:   data.New(db),
		mailer:   mailer.New(config.smtp.host, config.smtp.port, config.smtp.username, config.smtp.password, config.smtp.sender),
	}

	srv := http.Server{
		Addr:     fmt.Sprintf(":%d", config.port),
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	app.infoLog.Printf("server started on port %d", config.port)
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
