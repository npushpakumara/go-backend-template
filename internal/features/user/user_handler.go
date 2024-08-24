package user

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/npushpakumara/go-backend-template/internal/config"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService}
}

func UserRouter(configs *config.Config, router *gin.Engine, handler *UserHandler, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := router.Group("api/v1")

	v1.Use(authMiddleware.MiddlewareFunc())
	{
		v1.GET("/users", handler.getAllUsers)
	}

}

func (uh *UserHandler) register(ctx *gin.Context) {}
func (uh *UserHandler) getAllUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}
