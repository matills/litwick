package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/matills/litwick/internal/config"
)

type SupabaseClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// VerifySupabaseToken verifies the JWT token from Supabase and returns the user ID and email
func VerifySupabaseToken(tokenString string) (userID string, email string, err error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.SupabaseJWTSecret), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*SupabaseClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	return claims.Sub, claims.Email, nil
}

// ParseSupabaseUser parses user data from Supabase webhook/response
type SupabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func ParseSupabaseUser(data []byte) (*SupabaseUser, error) {
	var user SupabaseUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
