package store

import (
	"context"

	"github.com/romik1505/auth/internal/model"
)

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
