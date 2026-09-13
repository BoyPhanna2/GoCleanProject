package domain

import "context"

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type AuthUsecase interface {
	Register(ctx context.Context, email, password, name string) (*User, error)
	Login(ctx context.Context, email, password string) (string, error)
}
