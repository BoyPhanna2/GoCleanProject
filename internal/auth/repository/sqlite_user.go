package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"myapp/internal/auth/domain"
	sharedDomain "myapp/internal/shared/domain"
)

type sqliteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) domain.UserRepository {
	// Auto-create table
	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL
		);
	`)
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) Create(ctx context.Context, user *domain.User) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO users (email, password, name) VALUES (?, ?, ?)", user.Email, user.Password, user.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return sharedDomain.ErrConflict
		}
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *sqliteUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, email, password, name FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
