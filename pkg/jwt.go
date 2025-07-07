// Package pkg provides core utilities and services used throughout the Everato application.
// This includes JWT handling, logging, template utilities, and other shared functionality.
package pkg

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// TokenSigner provides functionality for signing and verifying JWT tokens.
// It encapsulates the signing method and secret key used for JWT operations.
type TokenSigner struct {
	Key    string                 // The secret key used for signing the JWT
	Method *jwt.SigningMethodHMAC // The signing method used for the JWT
}

// NewTokenSigner creates a new TokenSigner instance with the provided secret key.
// It uses HMAC-SHA256 as the default signing method for security and compatibility.
//
// Parameters:
//   - key: The secret key used to sign and verify JWT tokens
//
// Returns:
//   - A configured TokenSigner instance ready to sign and verify tokens
func NewTokenSigner(key string) *TokenSigner {
	return &TokenSigner{
		Key:    key,                    // Secret key provided by the caller
		Method: jwt.SigningMethodHS256, // HMAC-SHA256 symmetric algorithm for security and compatibility
	}
}

// Sign generates a signed JWT token with the provided payload.
// It uses the configured signing method and key to create a secure token.
//
// Parameters:
//   - payload: Map of claims to include in the JWT (e.g., user ID, expiration time)
//
// Returns:
//   - The signed JWT token as a string
//   - An error if signing fails
func (ts *TokenSigner) Sign(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(ts.Method, payload)
	return token.SignedString([]byte(ts.Key))
}

// Verify checks the validity of the provided JWT token and returns the claims if valid.
// It verifies the token signature, checks the signing method, and extracts the claims.
//
// This method performs the following validations:
// 1. Verifies that the token is properly signed with the expected method
// 2. Validates the token signature using the configured secret key
// 3. Extracts and returns the claims if the token is valid
//
// Parameters:
//   - tokenString: The JWT token string to verify
//
// Returns:
//   - The token claims as a jwt.MapClaims object if valid
//   - An error if the token is invalid, expired, or uses an unexpected signing method
func (ts *TokenSigner) Verify(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Verify that the token uses the expected signing method (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, fatal: %v", jwt.ErrSignatureInvalid)
		}
		return []byte(ts.Key), nil
	})

	// Check for parsing errors or invalid token
	if err != nil || !token.Valid {
		return nil, err
	}

	// Extract and return claims if present and correctly typed
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid claims, fatal: %v", jwt.ErrTokenInvalidClaims)
}
