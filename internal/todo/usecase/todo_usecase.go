package usecase

import (
	"context"

	"myapp/internal/todo/domain"
	sharedDomain "myapp/internal/shared/domain"
)

type todoUsecase struct {
	repo domain.TodoRepository
}

func NewTodoUsecase(repo domain.TodoRepository) domain.TodoUsecase {
	return &todoUsecase{repo: repo}
}

func (u *todoUsecase) Create(ctx context.Context, userID int64, title, description string) (*domain.Todo, error) {
	if title == "" {
		return nil, sharedDomain.ErrBadRequest
	}

	todo := &domain.Todo{
		UserID:      userID,
		Title:       title,
		Description: description,
		Done:        false,
	}

	if err := u.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (u *todoUsecase) Get(ctx context.Context, userID, id int64) (*domain.Todo, error) {
	todo, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if todo.UserID != userID {
		return nil, sharedDomain.ErrNotFound // Return 404 instead of 403/401 to not leak existence
	}

	return todo, nil
}

func (u *todoUsecase) List(ctx context.Context, userID int64) ([]*domain.Todo, error) {
	return u.repo.ListByUserID(ctx, userID)
}

func (u *todoUsecase) Update(ctx context.Context, userID, id int64, title, description string, done bool) (*domain.Todo, error) {
	todo, err := u.Get(ctx, userID, id)
	if err != nil {
		return nil, err // Returns ErrNotFound if not found or not owned
	}

	todo.Title = title
	todo.Description = description
	todo.Done = done

	if err := u.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (u *todoUsecase) Delete(ctx context.Context, userID, id int64) error {
	_, err := u.Get(ctx, userID, id)
	if err != nil {
		return err // Returns ErrNotFound if not found or not owned
	}

	return u.repo.Delete(ctx, id)
}
