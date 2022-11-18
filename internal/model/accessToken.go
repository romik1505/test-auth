package model

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AccessTokenTTL = time.Minute * 30
)

var (
	JwtPrivateKey = []byte("JWT_PRIVATE_KEY")
)

type JwtClaims struct {
	UserID         string `json:"user_id"`
	RefreshTokenID string `json:"refresh_token_id"`
	jwt.StandardClaims
}

func NewSignedAccessToken(ctx context.Context, userId, refreshTokenID string) (string, error) {
	claims := JwtClaims{
		UserID:         userId,
		RefreshTokenID: refreshTokenID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(JwtPrivateKey)
}

func JWTKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(JwtPrivateKey), nil
}
