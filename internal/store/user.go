package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/romik1505/auth/internal/model"
)

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
