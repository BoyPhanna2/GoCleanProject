package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"myapp/internal/config"
	sharedMiddleware "myapp/internal/shared/middleware"
)

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>API Documentation</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/swagger.yaml',
        dom_id: '#swagger-ui',
      });
    };
  </script>
</body>
</html>`

func NewGinEngine(logger *zap.Logger, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(sharedMiddleware.Logger(logger))
	r.Use(gin.Recovery())
	r.Use(sharedMiddleware.CORSMiddleware())
	r.Use(sharedMiddleware.ErrorHandler())

	// Serve the swagger YAML file dynamically with replaced BASE_URL
	r.GET("/swagger.yaml", func(c *gin.Context) {
		content, err := os.ReadFile("./swagger.yaml")
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to read swagger file")
			return
		}
		
		replaced := strings.ReplaceAll(string(content), "${BASE_URL}", cfg.BaseURL)
		c.Data(http.StatusOK, "application/yaml", []byte(replaced))
	})

	// Serve the swagger UI interface
	r.GET("/swagger", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
	})

	return r
}

func StartHTTPServer(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger, engine *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting http server", zap.String("port", cfg.ServerPort))
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping http server")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})
}

var Module = fx.Options(
	fx.Provide(NewGinEngine),
	fx.Invoke(StartHTTPServer),
)
