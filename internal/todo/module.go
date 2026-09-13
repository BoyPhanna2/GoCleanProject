package todo

import (
	"go.uber.org/fx"
	
	"myapp/internal/todo/delivery/http"
	"myapp/internal/todo/repository"
	"myapp/internal/todo/usecase"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewSQLiteTodoRepository,
		usecase.NewTodoUsecase,
	),
	fx.Invoke(http.NewTodoHandler),
)
