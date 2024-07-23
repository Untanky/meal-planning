package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

func connectDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "../data/meal-planner.db")
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("reason", err))
		panic(err)
	}

	return db
}

func migrateDatabase(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS meals (date TEXT PRIMARY KEY, breakfast TEXT, lunch TEXT, dinner TEXT, snacks TEXT)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS nutrition (date TEXT PRIMARY KEY, calories INT, weight INT)`)
	if err != nil {
		panic(err)
	}
}
