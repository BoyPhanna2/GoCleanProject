package domain

import "context"

type Todo struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TodoRepository interface {
	Create(ctx context.Context, todo *Todo) error
	GetByID(ctx context.Context, id int64) (*Todo, error)
	ListByUserID(ctx context.Context, userID int64) ([]*Todo, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id int64) error
}

type TodoUsecase interface {
	Create(ctx context.Context, userID int64, title, description string) (*Todo, error)
	Get(ctx context.Context, userID, id int64) (*Todo, error)
	List(ctx context.Context, userID int64) ([]*Todo, error)
	Update(ctx context.Context, userID, id int64, title, description string, done bool) (*Todo, error)
	Delete(ctx context.Context, userID, id int64) error
}
