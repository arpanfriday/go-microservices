package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "80"
const waitTillRetry = 2
const maxRetries = 10

var retries int64

type Config struct {
	Repo data.Repository
}

func main() {
	log.Println("Starting autentication service")

	// Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to Postgres")
	}

	// Setup config
	app := Config{}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("PostgreSQL not yet ready")
			retries++
		} else {
			log.Println("Connected to PostgreSQL")
			return connection
		}

		if retries > maxRetries {
			log.Println(err)
			return nil
		}

		log.Printf("Retrying in %d seconds ... ", waitTillRetry)
		time.Sleep(waitTillRetry * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
