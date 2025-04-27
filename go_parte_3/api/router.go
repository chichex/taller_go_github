package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"parte3/internal/user"
)

// InitRoutes registers all user CRUD endpoints on the given Gin engine.
// It initializes the storage, service, and handler, then binds each HTTP
// method and path to the appropriate handler function.
func InitRoutes(e *gin.Engine) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	storage := user.NewLocalStorage()
	service := user.NewService(storage, logger)

	h := handler{
		userService: service,
		logger:      logger,
	}

	e.POST("/users", h.handleCreate)
	e.GET("/users/:id", h.handleRead)
	e.PATCH("/users/:id", h.handleUpdate)
	e.DELETE("/users/:id", h.handleDelete)

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
