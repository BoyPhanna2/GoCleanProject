package repository

import (
	"context"
	"database/sql"
	"errors"

	"myapp/internal/todo/domain"
	sharedDomain "myapp/internal/shared/domain"
)

type sqliteTodoRepository struct {
	db *sql.DB
}

func NewSQLiteTodoRepository(db *sql.DB) domain.TodoRepository {
	// Auto-create table
	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`)
	return &sqliteTodoRepository{db: db}
}

func (r *sqliteTodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO todos (user_id, title, description, done) VALUES (?, ?, ?, ?)",
		todo.UserID, todo.Title, todo.Description, todo.Done)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	todo.ID = id
	return nil
}

func (r *sqliteTodoRepository) GetByID(ctx context.Context, id int64) (*domain.Todo, error) {
	var todo domain.Todo
	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, title, description, done FROM todos WHERE id = ?", id).
		Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (r *sqliteTodoRepository) ListByUserID(ctx context.Context, userID int64) ([]*domain.Todo, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, title, description, done FROM todos WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*domain.Todo
	for rows.Next() {
		var todo domain.Todo
		if err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Done); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (r *sqliteTodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	res, err := r.db.ExecContext(ctx, "UPDATE todos SET title = ?, description = ?, done = ? WHERE id = ?",
		todo.Title, todo.Description, todo.Done, todo.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sharedDomain.ErrNotFound
	}
	return nil
}

func (r *sqliteTodoRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sharedDomain.ErrNotFound
	}
	return nil
}
