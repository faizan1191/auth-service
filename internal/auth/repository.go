package auth

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash)
	return err
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `
	    SELECT id, email, password_hash, created_at
	    FROM users
	    WHERE email = $1
	`
	row := r.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no user found
		}
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByID(userID string) (*User, error) {
	user := &User{}

	query := `
	    SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdatePassword(userId, password string) error {
	query := `
		UPDATE users SET password_hash=$1
		WHERE id=$2
	`
	_, err := r.db.Exec(query, password, userId)
	return err
}
