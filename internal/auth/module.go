package auth

import (
	"go.uber.org/fx"
	
	"myapp/internal/auth/delivery/http"
	"myapp/internal/auth/repository"
	"myapp/internal/auth/usecase"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewSQLiteUserRepository,
		usecase.NewAuthUsecase,
	),
	fx.Invoke(http.NewAuthHandler),
)
