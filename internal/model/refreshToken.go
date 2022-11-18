package model

import (
	"database/sql"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	RefreshPrivateKey = []byte("REFRESH_PRIVATE_KEY")
)

type RefreshSession struct {
	ID           sql.NullString `db:"id"`
	UserID       sql.NullString `db:"user_id"`
	RefreshToken []byte         `db:"refresh_token"`
	ExpiresIn    sql.NullInt64  `db:"expires_in"`
	CreatedAT    sql.NullTime   `db:"created_at"`
}

type RefreshTokenClaims struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func (r *RefreshSession) GenerateRefreshTokenString() (string, error) {
	claims := RefreshTokenClaims{
		ID:     r.ID.String,
		UserID: r.UserID.String,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: r.ExpiresIn.Int64,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString(RefreshPrivateKey)
	if err != nil {
		return "", err
	}

	r.RefreshToken = []byte(signedToken)

	return signedToken, nil
}

func (r *RefreshSession) HashToken() error {
	hashedToken, err := bcrypt.GenerateFromPassword(r.RefreshToken, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	r.RefreshToken = hashedToken
	return nil
}
