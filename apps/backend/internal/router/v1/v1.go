package v1

import (
	"github.com/aneesh1213/tasker/internal/handler"
	"github.com/aneesh1213/tasker/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(router *echo.Group, handlers *handler.Handlers, middleware *middleware.Middlewares) {
	// Register todo routes
	registerTodoRoutes(router, handlers.Todo, handlers.Comment, middleware.Auth)

	// Register category routes
	registerCategoryRoutes(router, handlers.Category, middleware.Auth)

	// Register comment routes
	registerCommentRoutes(router, handlers.Comment, middleware.Auth)
}
