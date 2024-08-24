package middlewares

import (
	"time"

	"github.com/npushpakumara/go-backend-template/internal/features/auth"
	userDto "github.com/npushpakumara/go-backend-template/internal/features/user/dto"
	apiError "github.com/npushpakumara/go-backend-template/pkg/errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/npushpakumara/go-backend-template/internal/config"
	"github.com/npushpakumara/go-backend-template/internal/features/auth/dto"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

var identityKey = "id"

func NewAuthMiddleware(as auth.AuthService, cfg *config.Config) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(cfg.JWT.Secret),
		Timeout:     cfg.JWT.AccessTokenExpiry,
		MaxRefresh:  cfg.JWT.RefreshTokenExpiry,
		IdentityKey: identityKey,
		TokenLookup: "cookie:access_token",
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			logger := logging.FromContext(ctx)
			var requestBody dto.SignInRequestDto

			if err := ctx.ShouldBindJSON(&requestBody); err != nil {
				logger.Errorw("api.middlewares.AuthMiddleware failed to get request body: v", err)
				return nil, jwt.ErrMissingLoginValues
			}

			userID, err := as.LoginUser(ctx, &requestBody)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &userDto.UserResponseDto{ID: userID}, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, apiError.ErrorResponse{Status: "error", Message: message})
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*userDto.UserResponseDto); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &userDto.UserResponseDto{
				ID: claims[identityKey].(string),
			}
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*userDto.UserResponseDto); ok && v.ID != "" {
				return true
			}
			return false
		},
		LoginResponse: func(c *gin.Context, code int, token string, expires time.Time) {
			c.SetCookie("access_token", token, int(expires.Sub(time.Now()).Seconds()), "/", "", false, true)
			c.JSON(code, apiError.ErrorResponse{Status: "success", Message: "Login successfully"})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.SetCookie("access_token", "", -1, "/", "", false, true)
			c.JSON(code, apiError.ErrorResponse{Status: "success", Message: "Logout successfully"})
		},

		RefreshResponse: func(c *gin.Context, code int, token string, expires time.Time) {
			c.SetCookie("access_token", token, int(expires.Sub(time.Now()).Seconds()), "/", "", false, true) // Set as HTTP-only
			c.JSON(code, apiError.ErrorResponse{Status: "success", Message: "Token refresh successfully"})

		},
	})
}
