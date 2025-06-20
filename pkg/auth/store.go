package auth

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type AuthStore struct {
	DB *sql.DB
}

func (a *AuthStore) Migrate() error {
	_, err := a.DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password_hash TEXT,
		role TEXT
	);`)
	return err
}

func (a *AuthStore) CreateUser(u User) error {
	_, err := a.DB.Exec(`INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)`,
		u.Username, u.PasswordHash, u.Role)
	return err
}

func (a *AuthStore) FindUser(username string) (*User, error) {
	row := a.DB.QueryRow(`SELECT id, username, password_hash, role FROM users WHERE username = ?`, username)
	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role); err != nil {
		return nil, err
	}
	return &u, nil
}
