package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/npushpakumara/go-backend-template/internal/features/auth/dto"
	"github.com/npushpakumara/go-backend-template/internal/postgres"
	"github.com/npushpakumara/go-backend-template/pkg"
	apiError "github.com/npushpakumara/go-backend-template/pkg/errors"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService AuthService
}

// NewAuthHandler creates a new instance of AuthHandler with the given AuthService
func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// AuthRouter sets up the routes for authentication-related API endpoints
// It groups the routes under "api/v1/auth" and assigns handler functions to the routes
func AuthRouter(router *gin.Engine, handler *AuthHandler) {
	v1 := router.Group("api/v1/auth")

	v1.Use()
	{
		v1.POST("/sign-up", handler.signUpUser)
		v1.GET("/verify", handler.verifyUser)
	}
}

// signUpUser handles the user registration request
// It parses the JSON request body, validates it, and calls the authService to register the user
func (ah *AuthHandler) signUpUser(ctx *gin.Context) {
	logger := logging.FromContext(ctx)
	var requestBody dto.SignUpRequestDto

	// Bind and validate the JSON request body
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		logger.Errorw("auth.handler.signUpUser failed to get request body: v", err)
		var details []*pkg.ValidationErrDetail
		if vErrs, ok := err.(validator.ValidationErrors); ok {
			details = pkg.ValidationErrorDetails(&requestBody, "json", vErrs)
		}
		ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "error", Message: "Invalid request body", Errors: details})
		return
	}

	// Call the AuthService to register the user
	err := ah.authService.RegisterUser(ctx, &requestBody)
	if err != nil {
		if errors.Is(err, postgres.ErrKeyDuplicate) {
			ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "error", Message: "User already exist in the system", Errors: nil})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiError.ErrorResponse{Status: "error", Message: "Failed to signup user", Errors: nil})
		return
	}
}

// verifyUser handles the user verification request
// It extracts the token from the query parameters and calls the authService to activate the user's account
func (ah *AuthHandler) verifyUser(ctx *gin.Context) {
	logger := logging.FromContext(ctx)

	// Get the token from query parameters
	token, ok := ctx.GetQuery("token")
	if !ok {
		logger.Error("auth.handler.VerifyUser failed to get token")

		ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "Missing or invalid token", Errors: nil})
		return
	}

	// Call the AuthService to activate the account
	if err := ah.authService.ActivateAccount(ctx, token); err != nil {
		if errors.Is(err, postgres.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "User not found", Errors: nil})
			return
		}

		ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "Missing or invalid token", Errors: nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
