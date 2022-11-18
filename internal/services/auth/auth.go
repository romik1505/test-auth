package auth

import (
	"github.com/romik1505/auth/internal/store"
)

type AuthService struct {
	Storage store.Storage
}

func NewAuthService(s store.Storage) *AuthService {
	return &AuthService{
		Storage: s,
	}
}
