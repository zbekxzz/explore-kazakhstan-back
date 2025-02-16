package jwt

import (
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
	Role  string `json:"role"`
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
