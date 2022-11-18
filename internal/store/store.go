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

type IStorage interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUser(ctx context.Context, login string) (model.User, error)
	CreateSession(ctx context.Context, session model.RefreshSession) (model.RefreshSession, error)
}
