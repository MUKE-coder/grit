package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles JWT token operations.
type AuthService struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// Claims represents JWT claims.
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateTokenPair creates a new access + refresh token pair.
func (s *AuthService) GenerateTokenPair(userID uint, email, role string) (*TokenPair, error) {
	accessToken, expiresAt, err := s.generateToken(userID, email, role, s.AccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, _, err := s.generateToken(userID, email, role, s.RefreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// ValidateToken parses and validates a JWT token.
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GenerateResetToken creates a random hex token for password resets.
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generating reset token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) generateToken(userID uint, email, role string, expiry time.Duration) (string, int64, error) {
	expiresAt := time.Now().Add(expiry)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}
