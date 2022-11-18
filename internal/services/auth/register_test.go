package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/romik1505/auth/internal/mapper"
	"github.com/romik1505/auth/internal/model"
	"github.com/romik1505/auth/internal/store"
	mock_store "github.com/romik1505/auth/pkg/mock/store/mock_storage"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_store.NewMockIStorage(ctrl)
	authService := NewAuthService(s)

	testCases := []struct {
		name       string
		request    mapper.RegisterRequest
		wantErr    error
		hookBefore func()
	}{
		{
			name: "register new user",
			request: mapper.RegisterRequest{
				Login:    "user_login",
				Email:    "email@ya.ru",
				Password: "12345678",
				Phone:    "+71234567890",
			},
			hookBefore: func() {
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, user model.User) (model.User, error) {
					user.ID = store.NewNullString(uuid.NewString())
					return user, nil
				})
			},
			wantErr: nil,
		},
		{
			name: "user exist",
			request: mapper.RegisterRequest{
				Login:    "user_login",
				Email:    "email@ya.ru",
				Password: "12345678",
				Phone:    "+71234567890",
			},
			hookBefore: func() {
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(model.User{}, fmt.Errorf("some error"))
			},
			wantErr: ErrUserAlreadyExist,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.hookBefore()
			err := authService.Register(context.Background(), tt.request)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
