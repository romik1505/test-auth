package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/romik1505/auth/internal/mapper"
	"github.com/romik1505/auth/internal/model"
	"github.com/romik1505/auth/internal/store"
	"golang.org/x/crypto/bcrypt"
)

const (
	RefreshTokenTTL = time.Hour * 24 * 30 // 30 days
)

func (a AuthService) Login(ctx context.Context, req mapper.LoginRequest) (mapper.TokenPair, error) {
	u, err := a.Storage.GetUser(ctx, req.Login)
	if err != nil {
		return mapper.TokenPair{
			Status: "user not found",
		}, nil
	}

	log.Println(u)

	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return mapper.TokenPair{
				Status: "invalid password",
			}, nil
		}
	}

	refreshSession := model.RefreshSession{
		UserID:    u.ID,
		ExpiresIn: store.NewNullInt64(time.Now().Add(RefreshTokenTTL).Unix()),
	}
	refreshSession, err = a.Storage.CreateSession(ctx, refreshSession)
	if err != nil {
		return mapper.TokenPair{}, err
	}

	refreshToken, err := refreshSession.GenerateRefreshTokenString()
	if err != nil {
		return mapper.TokenPair{}, err
	}

	err = refreshSession.HashToken()
	if err != nil {
		return mapper.TokenPair{}, err
	}

	accessToken, err := model.NewSignedAccessToken(ctx, u.ID.String, refreshSession.ID.String)
	if err != nil {
		return mapper.TokenPair{}, err
	}
	log.Println("logged in")

	return mapper.TokenPair{
		Status:       "ok",
		AccessToken:  accessToken,
		RefreshToken: base64.StdEncoding.EncodeToString([]byte(refreshToken)),
	}, nil
}
