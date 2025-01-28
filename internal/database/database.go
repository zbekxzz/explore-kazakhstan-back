package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'Manager'
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
