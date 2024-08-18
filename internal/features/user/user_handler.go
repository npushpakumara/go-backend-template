package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/npushpakumara/go-backend-template/internal/config"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService}
}

func UserRouter(configs *config.Config, router *gin.Engine, handler *UserHandler) {
	v1 := router.Group("api/v1")

	v1.Use()
	{
		v1.POST("/users", handler.register)
		v1.GET("/users", handler.getAllUsers)
	}

}

func (uh *UserHandler) register(ctx *gin.Context) {}
func (uh *UserHandler) getAllUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}
