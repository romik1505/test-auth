package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/romik1505/auth/internal/model"
)

type Storage struct {
	DB *sqlx.DB
}

func (s Storage) Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(s.DB)
}

func (s Storage) GetUser(ctx context.Context, login string) (model.User, error) {
	query := s.Builder().Select("*").From("users").Where(sq.Eq{"login": login}).Limit(1)
	q, vars, err := query.ToSql()
	if err != nil {
		return model.User{}, err
	}

	user := model.User{}
	if err := s.DB.QueryRowxContext(ctx, q, vars...).StructScan(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s Storage) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	query := s.Builder().Insert("users").SetMap(map[string]interface{}{
		"login":    user.Login,
		"email":    user.Email,
		"password": user.Password,
		"phone":    user.Phone,
	}).Suffix("RETURNING id")
	q, vars, err := query.ToSql()
	if err != nil {
		return model.User{}, err
	}

	if err := s.DB.QueryRowxContext(ctx, q, vars...).StructScan(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s Storage) CreateSession(ctx context.Context, session model.RefreshSession) (model.RefreshSession, error) {
	query := s.Builder().Insert("sessions").SetMap(map[string]interface{}{
		"user_id":       session.UserID,
		"refresh_token": session.RefreshToken,
		"expires_in":    session.ExpiresIn,
	}).Suffix("RETURNING id, created_at")

	q, vars, err := query.ToSql()
	if err != nil {
		return model.RefreshSession{}, err
	}

	if err := s.DB.QueryRowxContext(ctx, q, vars...).StructScan(&session); err != nil {
		return model.RefreshSession{}, err
	}
	return session, nil
}
