package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"myapp/internal/todo/domain"
	sharedDomain "myapp/internal/shared/domain"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoRepository) GetByID(ctx context.Context, id int64) (*domain.Todo, error) {
	args := m.Called(ctx, id)
	if todo, ok := args.Get(0).(*domain.Todo); ok {
		return todo, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTodoRepository) ListByUserID(ctx context.Context, userID int64) ([]*domain.Todo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGet_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	uc := NewTodoUsecase(mockRepo)

	expectedTodo := &domain.Todo{ID: 1, UserID: 42, Title: "Test"}
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedTodo, nil)

	todo, err := uc.Get(context.Background(), 42, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodo, todo)
}

func TestGet_WrongUser(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	uc := NewTodoUsecase(mockRepo)

	expectedTodo := &domain.Todo{ID: 1, UserID: 42, Title: "Test"}
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedTodo, nil)

	todo, err := uc.Get(context.Background(), 99, 1) // Trying to access as user 99

	assert.ErrorIs(t, err, sharedDomain.ErrNotFound)
	assert.Nil(t, todo)
}
