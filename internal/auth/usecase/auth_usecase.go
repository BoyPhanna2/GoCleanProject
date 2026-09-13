package usecase

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"myapp/internal/auth/domain"
	"myapp/internal/config"
	sharedDomain "myapp/internal/shared/domain"
)

type authUsecase struct {
	repo domain.UserRepository
	cfg  *config.Config
}

func NewAuthUsecase(repo domain.UserRepository, cfg *config.Config) domain.AuthUsecase {
	return &authUsecase{
		repo: repo,
		cfg:  cfg,
	}
}

func (u *authUsecase) Register(ctx context.Context, email, password, name string) (*domain.User, error) {
	// Validate email
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("%w: invalid email format", sharedDomain.ErrBadRequest)
	}

	// Validate password length
	if len(password) < 8 {
		return nil, fmt.Errorf("%w: password must be at least 8 characters", sharedDomain.ErrBadRequest)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to hash password", sharedDomain.ErrInternal)
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == sharedDomain.ErrNotFound {
			return "", sharedDomain.ErrUnauthorized // Don't leak existence
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", sharedDomain.ErrUnauthorized
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("%w: failed to sign token", sharedDomain.ErrInternal)
	}

	return tokenString, nil
}
