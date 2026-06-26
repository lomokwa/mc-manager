package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	return migrate()
}

func migrate() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username VARCHAR(50) UNIQUE NOT NULL,
    	email VARCHAR(255) UNIQUE,
    	password_hash VARCHAR(255) NOT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS invitations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token VARCHAR(64) UNIQUE NOT NULL,
			email VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL,
			used_at TIMESTAMP
		);
	`)

	return err
}
