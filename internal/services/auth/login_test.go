package auth

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/romik1505/auth/internal/mapper"
	"github.com/romik1505/auth/internal/model"
	"github.com/romik1505/auth/internal/store"
	mock_store "github.com/romik1505/auth/pkg/mock/store/mock_storage"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_store.NewMockIStorage(ctrl)
	service := NewAuthService(s)

	testCases := []struct {
		name       string
		req        mapper.LoginRequest
		wantStatus string
		wantErr    error
		hookBefore func()
	}{
		{
			name: "ok case",
			req: mapper.LoginRequest{
				Login:    "login",
				Password: "12345678",
			},
			wantStatus: "ok",
			wantErr:    nil,
			hookBefore: func() {
				s.EXPECT().GetUser(gomock.Any(), "login").Return(model.User{
					ID:       store.NewNullString(uuid.NewString()),
					Login:    store.NewNullString("login"),
					Password: []byte("$2a$10$GBmeEuVNAyrHVSzuEc0vpOcu8p0xg8WVMKPQJ5OGC2oD.ByYYQCI2"), // valid hash for password
					Email:    store.NewNullString("email@ya.ru"),
					Phone:    store.NewNullString("+71234567890"),
				}, nil)

				s.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(model.RefreshSession{
					ID:           store.NewNullString(uuid.NewString()),
					UserID:       store.NewNullString(uuid.NewString()),
					RefreshToken: []byte{},
					ExpiresIn:    store.NewNullInt64(time.Now().Add(time.Hour * 900).Unix()),
					CreatedAt:    store.NewNullTime(time.Date(2022, 12, 10, 9, 10, 10, 44, time.UTC)),
				}, nil)
			},
		},
		{
			name: "user not found",
			req: mapper.LoginRequest{
				Login:    "login",
				Password: "12345678",
			},
			wantStatus: "user not found",
			wantErr:    ErrUserNotFound,
			hookBefore: func() {
				s.EXPECT().GetUser(gomock.Any(), "login").Return(model.User{}, sql.ErrNoRows)
			},
		},
		{
			name: "invalid password",
			req: mapper.LoginRequest{
				Login:    "login",
				Password: "12345678",
			},
			wantStatus: "invalid password",
			wantErr:    ErrInvalidPassword,
			hookBefore: func() {
				s.EXPECT().GetUser(gomock.Any(), "login").Return(model.User{
					ID:       store.NewNullString(uuid.NewString()),
					Login:    store.NewNullString("login"),
					Password: []byte("$2a$10$GBmeEuVNAyrHVSzuEc0vpOcu8p0xg8WVMKPQJ5OGC2oD.ByYYQAAA"), // invalid hash for password
					Email:    store.NewNullString("email@ya.ru"),
					Phone:    store.NewNullString("+71234567890"),
				}, nil)
			},
		},
		{
			name: "creation session error",
			req: mapper.LoginRequest{
				Login:    "login",
				Password: "12345678",
			},
			wantStatus: "",
			wantErr:    fmt.Errorf("some db error"),
			hookBefore: func() {
				s.EXPECT().GetUser(gomock.Any(), "login").Return(model.User{
					ID:       store.NewNullString(uuid.NewString()),
					Login:    store.NewNullString("login"),
					Password: []byte("$2a$10$GBmeEuVNAyrHVSzuEc0vpOcu8p0xg8WVMKPQJ5OGC2oD.ByYYQCI2"), // valid hash for password
					Email:    store.NewNullString("email@ya.ru"),
					Phone:    store.NewNullString("+71234567890"),
				}, nil)

				s.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(model.RefreshSession{}, fmt.Errorf("some db error"))
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.hookBefore()

			pair, err := service.Login(context.Background(), tt.req)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.wantStatus, pair.Status)
		})
	}
}
