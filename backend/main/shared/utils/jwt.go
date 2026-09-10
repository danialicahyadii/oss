package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID   uint64   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	JTI      string   `json:"jti"`

	jwt.RegisteredClaims
}

func GenerateToken(
	userID uint64,
	username string,
	roles []string,
) (string, error) {

	hours := 24

	if value := os.Getenv("JWT_EXPIRE_HOURS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			hours = parsed
		}
	}

	jti := uuid.New().String()
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: jti,
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(hours) * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}

func ParseToken(tokenString string) (*JWTClaims, error) {

	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
