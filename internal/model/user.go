package model

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       sql.NullString `db:"id"`
	Login    sql.NullString `db:"login"`
	Password []byte         `db:"password"`
	Email    sql.NullString `db:"email"`
	Phone    sql.NullString `db:"phone"`
}

func (u *User) HashPassword() error {
	hashedPass, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = hashedPass
	return nil
}
