package pkg

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type TokenSigner struct {
	Key    string                 // The secret key used for signing the JWT
	Method *jwt.SigningMethodHMAC // The signing method used for the JWT
}

// Instantiate a new TokenSigner instance
func NewTokenSigner(key string) *TokenSigner {
	return &TokenSigner{
		Key:    key,                    // Assuming that the key would be provided by the caller
		Method: jwt.SigningMethodHS256, // Symmetric algorithm for sake of simplicity
	}
}

// Sign generates a signed JWT token with the provided payload.
func (ts *TokenSigner) Sign(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(ts.Method, payload)
	return token.SignedString([]byte(ts.Key))
}

// Verify checks the validity of the provided JWT token and returns the claims if valid.
func (ts *TokenSigner) Verify(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, fatal: %v", jwt.ErrSignatureInvalid)
		}
		return []byte(ts.Key), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid claims, fatal: %v", jwt.ErrTokenInvalidClaims)
}
