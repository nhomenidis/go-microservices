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
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

// Config holds the configuration for the API
type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting Authentication service")

	// TODO connect to database
	conn := connectToDb()

	if conn == nil {
		log.Panic("Could not connect to database!")
	}
	// setup configuration
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// let's connect to the database
func openDb(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}

	// check if we can connect to the database
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDb() *sql.DB {
	databaseUrl := os.Getenv("DATABASE_URL")

	for {
		connection, err := openDb(databaseUrl)
		if err != nil {
			log.Println("Error connecting to database", err)
			time.Sleep(2 * time.Second)
			counts++
			if counts > 10 {
				return nil
			}
			continue
		} else {
			log.Println("Connected to database")
			return connection
		}
	}
}
