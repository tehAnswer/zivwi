package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var (
	DatabaseUrl = os.Getenv("DB_URL")
)

type Database struct {
	Url        string
	Connection *sql.DB
}

func NewDatabase() *Database {
	if connection, err := sql.Open("postgres", DatabaseUrl); err == nil {
		return &Database{
			Url:        DatabaseUrl,
			Connection: connection,
		}
	} else {
		panic("Unable to connect to database.\n" + err.Error())
	}
}
