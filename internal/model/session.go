package model

import "database/sql"

type Session struct {
	ID           sql.NullString `db:"id"`
	UserID       sql.NullString `db:"user_id"`
	RefreshToken []byte         `db:"refresh_token"`
	ExpiresIn    sql.NullInt64  `db:"expires_in"`
	CreatedAT    sql.NullTime   `db:"created_at"`
}
