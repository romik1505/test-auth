package mapper

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (l *LoginRequest) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Login, validation.Required),
		validation.Field(&l.Password, validation.Length(5, 50)),
	)
}

type TokenPair struct {
	Status       string `json:"status,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RegisterRequest struct {
	Login    string `json:"login,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

func (r *RegisterRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Login, validation.Required, validation.Length(5, 50)),
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Password, validation.Required, validation.Length(5, 50)),
		validation.Field(&r.Phone, validation.Required, validation.Length(5, 20)),
	)
}
