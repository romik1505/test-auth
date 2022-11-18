package store

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/romik1505/auth/internal/model"
	"github.com/stretchr/testify/require"
)

func TestStorage_CreateSession(t *testing.T) {
	mustTruncateAll()

	user, _ := storage.CreateUser(context.Background(), model.User{
		Login:    NewNullString("login"),
		Password: []byte("password"),
		Email:    NewNullString("ya@mail.ru"),
		Phone:    NewNullString("+78881112233"),
	})

	tests := []struct {
		name       string
		input      model.RefreshSession
		want       model.RefreshSession
		wantErr    error
		hookBefore func()
	}{
		{
			name: "first user session",
			input: model.RefreshSession{
				UserID:       user.ID,
				RefreshToken: []byte("token"),
				ExpiresIn:    NewNullInt64(time.Now().Unix()),
			},
		},
		{
			name: "second user session",
			input: model.RefreshSession{
				UserID:       user.ID,
				RefreshToken: []byte("token"),
				ExpiresIn:    NewNullInt64(time.Now().Unix()),
			},
		},
		{
			name: "create session with invalid user_id",
			input: model.RefreshSession{
				UserID:       user.ID,
				RefreshToken: []byte("token"),
				ExpiresIn:    NewNullInt64(time.Now().Unix()),
			},
			hookBefore: func() {
				mustTruncateAll()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.CreateSession(context.Background(), tt.input)
			require.Equal(t, tt.wantErr, err)
			require.Empty(t, cmp.Diff(tt.input, got, cmpopts.IgnoreFields(model.RefreshSession{}, "ID", "CreatedAt")))
			require.True(t, got.ID.Valid)
			require.True(t, got.CreatedAt.Valid)

		})
	}
	mustTruncateAll()
}
