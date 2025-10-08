package repository

import (
	"database/sql"
	"tugas-praktikum-crud/app/model"
)

func FindUserByUsername(db *sql.DB, username string) (*model.User, string, error) {
	var user model.User
	var passwordHash string

	query := "SELECT id, username, email, password_hash, role, created_at FROM users WHERE username = $1"

	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}