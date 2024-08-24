package user

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/npushpakumara/go-backend-template/internal/config"
)

// Handler struct represents the HTTP handler for user-related operations.
// It contains a reference to the userService which handles the business logic.
type Handler struct {
	userService Service
}

// NewUserHandler creates a new Handler instance with the provided userService.
// It returns a pointer to the Handler struct.
func NewUserHandler(userService Service) *Handler {
	return &Handler{userService}
}

// Router sets up the routes for the user-related API endpoints.
// It takes in the application configuration, the Gin router instance, the handler for user operations,
// and the authentication middleware to secure the endpoints.
func Router(configs *config.Config, router *gin.Engine, handler *Handler, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := router.Group("api/v1")

	v1.Use(authMiddleware.MiddlewareFunc())
	{
		v1.GET("/users", handler.getAllUsers)
	}

}

// getAllUsers is a handler method for the Handler struct.
func (uh *Handler) getAllUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}
