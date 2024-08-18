package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/npushpakumara/go-backend-template/pkg/errors"
)

// NewJwtToken creates a new JWT token with the given user ID, secret key, and expiration duration.
// It sets the issuer to "example.com", the subject to the provided user ID, and includes both issued and expiration dates in the token claims.
// The token is signed using the HS256 algorithm and the provided secret key.
// Returns the signed token string and an error if any occurred during signing.
func NewJwtToken(id, secret string, exp time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "example.com",
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ExtractSubjectFromToken parses the JWT token using the provided secret key to verify its validity.
// It ensures the token is signed with the HMAC signing method and extracts the "sub" (subject) claim from the token's claims.
// Returns the subject as a string and an error if the token is invalid or if any other error occurs during parsing.
func ExtractSubjectFromToken(secret, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is signed with the expected signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	// Assert the token claims to jwt.MapClaims type
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.ErrInvalidToken
	}

	// Extract the "sub" (subject) claim from the claims
	subject, ok := claims["sub"].(string)
	if !ok {
		return "", errors.ErrInvalidToken
	}

	return subject, nil
}
