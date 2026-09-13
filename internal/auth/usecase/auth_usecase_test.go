package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"myapp/internal/auth/domain"
	"myapp/internal/config"
	sharedDomain "myapp/internal/shared/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "secret"}
	uc := NewAuthUsecase(mockRepo, cfg)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := uc.Register(context.Background(), "test@example.com", "password123", "Test User")
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.NotEqual(t, "password123", user.Password) // Should be hashed
}

func TestRegister_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "secret"}
	uc := NewAuthUsecase(mockRepo, cfg)

	user, err := uc.Register(context.Background(), "invalid-email", "password123", "Test User")
	
	assert.ErrorIs(t, err, sharedDomain.ErrBadRequest)
	assert.Nil(t, user)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "secret"}
	uc := NewAuthUsecase(mockRepo, cfg)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	storedUser := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashed),
		Name:     "Test User",
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(storedUser, nil)

	token, err := uc.Login(context.Background(), "test@example.com", "password123")
	
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_Unauthorized(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "secret"}
	uc := NewAuthUsecase(mockRepo, cfg)

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, sharedDomain.ErrNotFound)

	token, err := uc.Login(context.Background(), "test@example.com", "password123")
	
	assert.ErrorIs(t, err, sharedDomain.ErrUnauthorized)
	assert.Empty(t, token)
}
