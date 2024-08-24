package auth

import (
	"errors"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/npushpakumara/go-backend-template/internal/config"
	"github.com/npushpakumara/go-backend-template/internal/features/auth/dto"
	"github.com/npushpakumara/go-backend-template/internal/postgres"
	"github.com/npushpakumara/go-backend-template/pkg"
	apiError "github.com/npushpakumara/go-backend-template/pkg/errors"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// Handler handles authentication-related requests
type Handler struct {
	authService Service
	cfg         *config.Config // Configuration settings for the application
}

// NewAuthHandler creates a new instance of Handler with the given Service
func NewAuthHandler(authService Service, cfg *config.Config) *Handler {
	return &Handler{authService, cfg}
}

// Router sets up the routes for authentication-related API endpoints
// It groups the routes under "api/v1/auth" and assigns handler functions to the routes
func Router(router *gin.Engine, handler *Handler, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := router.Group("api/v1")

	v1.Use()
	{
		// User authentication and management
		v1.POST("/auth/sign-up", handler.signUp)
		v1.POST("/auth/sign-in", authMiddleware.LoginHandler)
		v1.POST("/auth/sign-out", authMiddleware.LogoutHandler)
		v1.POST("/auth/refresh-token", authMiddleware.RefreshHandler)

		// Account verification and email management
		v1.GET("/auth/verify-email", handler.verifyUser)
		v1.POST("/auth/resend-verification-email", handler.reSendVerificationEmail)

		// Password management
		v1.PUT("/auth/reset-password", handler.resetPassword)

		// OAuth handling
		v1.GET("/oauth/:provider", OAuthMiddleware())
		v1.GET("/oauth/:provider/callback", OAuthCallbackMiddleware(authMiddleware, handler.authService.HandleOAuthUser))
	}

}

// signUpUser handles the user registration request
// It parses the JSON request body, validates it, and calls the authService to register the user
func (ah *Handler) signUp(ctx *gin.Context) {
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

	// Call the Service to register the user
	err := ah.authService.RegisterUser(ctx, &requestBody)
	if err != nil {
		if errors.Is(err, postgres.ErrKeyDuplicate) {
			ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "error", Message: "User already exist in the system", Errors: nil})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiError.ErrorResponse{Status: "error", Message: "Internal srver error", Errors: nil})
		return
	}
	ctx.JSON(http.StatusCreated, dto.SignUpResponseDto{Status: "success", Message: "User has been registered. Please check email for account confirmation"})
}

// verifyUser handles the user verification request
// It extracts the token from the query parameters and calls the authService to activate the user's account
func (ah *Handler) verifyUser(ctx *gin.Context) {
	logger := logging.FromContext(ctx)

	// Get the token from query parameters
	token, ok := ctx.GetQuery("token")
	if !ok {
		logger.Error("auth.handler.VerifyUser failed to get token")

		ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "Missing or invalid token", Errors: nil})
		return
	}

	// Call the Service to activate the account
	id, err := ah.authService.ActivateAccount(ctx, token)
	if err != nil {
		if errors.Is(err, postgres.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "User not found", Errors: nil})
			return
		}
		ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "Missing or invalid token", Errors: nil})
		return
	}

	if id == "" {
		logger.Error("auth.handler.VerifyUser failed to get user id")
		ctx.JSON(http.StatusInternalServerError, apiError.ErrorResponse{Status: "failed", Message: "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SignUpResponseDto{Status: "success", Message: "Account activated"})
}

// reSendVerificationEmail handles the request to resend the account verification email to the user.
// It expects the user's ID to be provided as a query parameter and performs the following steps:
func (ah *Handler) reSendVerificationEmail(ctx *gin.Context) {
	userID, ok := ctx.GetQuery("id")
	if !ok {
		ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "Missing user id", Errors: nil})
		return
	}

	user, err := ah.authService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "User not found", Errors: nil})
			return
		}
		ctx.JSON(http.StatusInternalServerError, apiError.ErrorResponse{Status: "failed", Message: "Internal server error", Errors: nil})
		return
	}

	if user.IsActive {
		ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "failed", Message: "User is already active", Errors: nil})
		return
	}

	err = ah.authService.SendAccountVerificationEmail(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiError.ErrorResponse{Status: "failed", Message: "Internal server error", Errors: nil})
		return
	}

	ctx.JSON(http.StatusOK, dto.SignUpResponseDto{Status: "success", Message: "Email has been sent"})
}

// resetPassword handles the request to reset a user's password.
// It expects a JSON body containing the user's current password and the new password.
func (ah *Handler) resetPassword(ctx *gin.Context) {
	logger := logging.FromContext(ctx)
	var requestBody dto.PasswordResetRequestDto

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		logger.Errorw("auth.handler.resetPassword failed to get request body: v", err)
		var details []*pkg.ValidationErrDetail
		if vErrs, ok := err.(validator.ValidationErrors); ok {
			details = pkg.ValidationErrorDetails(&requestBody, "json", vErrs)
		}
		ctx.JSON(http.StatusBadRequest, apiError.ErrorResponse{Status: "error", Message: "Invalid request body", Errors: details})
		return
	}

	err := ah.authService.ResetPassword(ctx, &requestBody)
	if err != nil {
		if errors.Is(err, apiError.ErrIncorrectPassword) {
			ctx.JSON(http.StatusUnauthorized, apiError.ErrorResponse{Status: "error", Message: "Invalid current password"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, apiError.ErrorResponse{Status: "failed", Message: "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SignUpResponseDto{Status: "success", Message: "Password updated successfully"})
}
