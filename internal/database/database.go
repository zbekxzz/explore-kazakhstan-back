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

	// Создание таблицы пользователей
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	// Создание таблицы событий (events)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			date DATE NOT NULL,
			time TIME NOT NULL,
			venue TEXT NOT NULL,
			description TEXT,
			note TEXT,
			price TEXT,
			image_url TEXT,
			attendees TEXT[], -- Хранит массив идентификаторов участников
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create events table:", err)
	}

	return db
}
