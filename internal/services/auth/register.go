package auth

import (
	"context"
	"errors"
	"log"

	"github.com/romik1505/auth/internal/mapper"
	"github.com/romik1505/auth/internal/model"
	"github.com/romik1505/auth/internal/store"
)

func (a AuthService) Register(ctx context.Context, req mapper.RegisterRequest) error {
	user := model.User{
		Login:    store.NewNullString(req.Login),
		Email:    store.NewNullString(req.Email),
		Password: []byte(req.Password),
		Phone:    store.NewNullString(req.Phone),
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	u, err := a.Storage.CreateUser(ctx, user)
	if err != nil {
		return ErrUserAlreadyExist
	}
	log.Println(u)
	return nil
}

var (
	ErrUserAlreadyExist = errors.New("user already exist")
)
