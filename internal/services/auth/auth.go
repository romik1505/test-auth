package auth

import (
	"context"
	"errors"
	"log"

	"github.com/romik1505/auth/internal/mapper"
	"github.com/romik1505/auth/internal/model"
	"github.com/romik1505/auth/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Storage store.Storage
}

func NewAuthService(s store.Storage) *AuthService {
	return &AuthService{
		Storage: s,
	}
}

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

	// TODO: create tokens
	log.Println("logged in")

	return mapper.TokenPair{
		Status: "ok",
	}, nil
}

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
