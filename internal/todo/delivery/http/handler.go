package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"myapp/internal/config"
	sharedMiddleware "myapp/internal/shared/middleware"
	"myapp/internal/shared/response"
	"myapp/internal/todo/domain"
)

type TodoHandler struct {
	usecase domain.TodoUsecase
}

func NewTodoHandler(r *gin.Engine, usecase domain.TodoUsecase, cfg *config.Config) {
	handler := &TodoHandler{usecase: usecase}

	todoGroup := r.Group("/api/v1/todos")
	todoGroup.Use(sharedMiddleware.JWTAuth(cfg))
	{
		todoGroup.POST("", handler.Create)
		todoGroup.GET("", handler.List)
		todoGroup.GET("/:id", handler.Get)
		todoGroup.PUT("/:id", handler.Update)
		todoGroup.DELETE("/:id", handler.Delete)
	}
}

func getUserID(c *gin.Context) int64 {
	return c.GetInt64("user_id")
}

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (h *TodoHandler) Create(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	todo, err := h.usecase.Create(c.Request.Context(), getUserID(c), req.Title, req.Description)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, todo)
}

func (h *TodoHandler) List(c *gin.Context) {
	todos, err := h.usecase.List(c.Request.Context(), getUserID(c))
	if err != nil {
		c.Error(err)
		return
	}

	// Handle nil slice for empty list
	if todos == nil {
		todos = []*domain.Todo{}
	}

	response.Success(c, http.StatusOK, todos)
}

func (h *TodoHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", "")
		c.Abort()
		return
	}

	todo, err := h.usecase.Get(c.Request.Context(), getUserID(c), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, todo)
}

type UpdateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", "")
		c.Abort()
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	todo, err := h.usecase.Update(c.Request.Context(), getUserID(c), id, req.Title, req.Description, req.Done)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", "")
		c.Abort()
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), getUserID(c), id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, nil)
}
