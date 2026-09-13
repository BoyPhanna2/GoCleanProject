package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"myapp/internal/auth"
	"myapp/internal/config"
	"myapp/internal/server"
	"myapp/internal/todo"
	"myapp/pkg/database"
	"myapp/pkg/logger"
)

func main() {
	app := fx.New(
		config.Module,
		logger.Module,
		database.Module,
		auth.Module,
		todo.Module,
		server.Module,
		fx.Invoke(func(l *zap.Logger) {
			l.Info("Application initialized successfully")
		}),
	)
	
	app.Run()
}
