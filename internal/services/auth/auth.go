package auth

import (
	"github.com/romik1505/auth/internal/store"
)

type AuthService struct {
	Storage store.IStorage
}

func NewAuthService(s store.IStorage) *AuthService {
	return &AuthService{
		Storage: s,
	}
}
