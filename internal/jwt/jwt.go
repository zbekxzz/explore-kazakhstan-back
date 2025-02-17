package jwt

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey               string
	accessTokenExpiresTime  time.Duration
	refreshTokenExpiresTime time.Duration
)

func InitJWT(key string, AccessTime time.Duration, refreshTime time.Duration) {
	secretKey = key
	accessTokenExpiresTime = AccessTime
	refreshTokenExpiresTime = refreshTime
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateTokens(email string) (string, string, error) {
	accessTokenClaims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExpiresTime).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(refreshTokenExpiresTime).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}

func ExtractUserEmail(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	tokenString := parts[1]

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	return claims.Email, nil
}
