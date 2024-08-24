package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/npushpakumara/go-backend-template/internal/features/auth/dto"
	"github.com/npushpakumara/go-backend-template/pkg/errors"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
)

// OAuthMiddleware is a Gin middleware function that handles the initial OAuth request.
// It sets up the necessary state for the OAuth flow, including setting the provider in the context
// and generating a state cookie to prevent CSRF attacks.
func OAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		if provider == "" {
			c.JSON(http.StatusBadRequest, errors.ErrorResponse{Status: "error", Message: "Provider not specified"})
			return
		}

		// Generate a random state string for the OAuth flow to prevent CSRF attacks.
		state := generateStateOauthCookie()
		// Set the state as a secure, HttpOnly cookie that expires in 5 minutes.
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Expires:  time.Now().Add(5 * time.Minute),
			HttpOnly: true,
			Secure:   true,
		})

		// Add the state parameter to the URL query string for the OAuth request.
		q := c.Request.URL.Query()
		q.Add("state", state)
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()

		gothic.BeginAuthHandler(c.Writer, c.Request)
		c.Abort()
	}
}

// OAuthCallbackMiddleware is a Gin middleware function that handles the callback from the OAuth provider.
// It completes the OAuth authentication, validates the state, and generates a JWT for the authenticated user.
func OAuthCallbackMiddleware(authMiddleware *jwt.GinJWTMiddleware, handleUser func(ctx context.Context, user goth.User) (*dto.OAuthResponseDto, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve a logger from the request context for logging purposes.
		logger := logging.FromContext(c.Request.Context())

		// Complete the OAuth authentication and retrieve the user information from the provider.
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			logger.Errorf("auth.middlewares.OAuthCallbackMiddleware failed to authenticate: %v", err.Error())
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{Status: "error", Message: "Authentication failed"})
			return
		}

		// Retrieve the state cookie from the request.
		cookie, err := c.Cookie("oauth_state")
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.ErrorResponse{Status: "error", Message: "State cookie not found"})
			return
		}

		// Validate the state parameter from the URL against the state stored in the cookie.
		state := c.Query("state")
		if state == "" || state != cookie {
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{Status: "error", Message: "Invalid state"})
			return
		}

		// Handle the authenticated user by invoking the provided handler function.
		result, err := handleUser(c.Request.Context(), user)
		if err != nil { // Handle any errors that occur during user handling.
			logger.Error("auth.middlewares.OAuthCallbackMiddleware failed to handle user", "error", err.Error())
			c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Status: "error", Message: "Internal server error"})
			return
		}

		// Generate a JWT token for the authenticated user using the provided JWT middleware.
		token, expires, err := authMiddleware.TokenGenerator(result.ID)
		if err != nil {
			logger.Error("auth.middlewares.OAuthCallbackMiddleware failed to handle user", "error", err.Error())
			c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Status: "error", Message: "Internal server error"})
			return
		}

		c.SetCookie("access_token", token, int(time.Until(expires).Seconds()), "/", "", false, true)

		c.JSON(http.StatusOK, dto.SignUpResponseDto{Status: "success", Message: "Successfully signed in"})
	}
}

// generateStateOauthCookie generates a random state string to be used in the OAuth flow.
// This state string is encoded in base64 and is used to protect against CSRF attacks.
func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
