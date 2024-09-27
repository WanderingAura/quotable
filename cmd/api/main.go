// The main package for quotable. Contains the code for starting the application, serving http requests
// and handling endpoints.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/WanderingAura/quotable/internal/data"
	"github.com/WanderingAura/quotable/internal/mailer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// The value is burnt into the executable at compile time.
var (
	buildTime string
	version   string
)

// Stores all the relevant info about the app (to be used by the handlers)
type application struct {
	config config
	logger *zerolog.Logger
	models data.Models    // Exposes CRUD operations on database tables
	mailer mailer.Mailer  // Used for sending an email after user registration
	wg     sync.WaitGroup // Used for graceful shutdown
}

// Used to configure the various settings of the app on start up
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

// Initialises the application and sets it up to start listening for HTTP reqs
func main() {
	var config config

	flag.IntVar(&config.port, "port", 4000, "API server port")

	flag.StringVar(&config.env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&config.debug, "debug", false, "debug mode")

	flag.StringVar(&config.logPath, "log-path", "./logs/quotable.log", "File to write error logs in")

	// database configs
	flag.StringVar(&config.db.dsn, "db-dsn", "", "Postgres database source name")
	flag.IntVar(&config.db.maxOpenConnections, "db-max-open-conns", 25, "Postgres max open connections")
	flag.IntVar(&config.db.maxIdleConnections, "db-max-idle-conns", 25, "Postgres max idle connections")
	flag.StringVar(&config.db.maxIdleDuration, "db-max-idle-time", "15m", "Postgres max connection idle time")

	// rate limiting config
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

	logFile := os.Stderr
	if !config.debug {
		logFile, err := os.Open(config.logPath)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(logFile).With().Timestamp().Logger()

	db, err := openDB(config)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("")
	}
	logger.Info().Msg("database connection successful!")

	app := &application{
		config: config,
		logger: &logger,
		models: data.New(db),
		mailer: mailer.New(config.smtp.host, config.smtp.port, config.smtp.username, config.smtp.password, config.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}

// Sets up the postgres database connection
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
