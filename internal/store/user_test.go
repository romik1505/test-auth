package store

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lib/pq"
	"github.com/romik1505/auth/internal/model"
	"github.com/stretchr/testify/require"
)

func TestStorage_CreateUser(t *testing.T) {
	mustTruncateAll()

	testCases := []struct {
		name       string
		input      model.User
		hookBefore func()
	}{
		{
			name: "create new user",
			input: model.User{
				Login:    NewNullString("login"),
				Password: []byte("hashed password"),
				Email:    NewNullString("email2@ya.ru"),
				Phone:    NewNullString("+79991234567"),
			},
			hookBefore: func() {
			},
		},
		{
			name: "try dublicate user login",
			input: model.User{
				Login:    NewNullString("login"),
				Password: []byte("hashed password"),
				Email:    NewNullString("email@ya.ru"),
				Phone:    NewNullString("+79991234567"),
			},
			hookBefore: func() {
				_, err := storage.CreateUser(context.Background(), model.User{
					Login:    NewNullString("login"),
					Password: []byte("hashed password"),
					Email:    NewNullString("email2@ya.ru"),
					Phone:    NewNullString("+79991234567"),
				})
				require.Nil(t, err)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.hookBefore()

			resp, err := storage.CreateUser(context.Background(), tt.input)
			if err != nil {
				require.Equal(t, err.(*pq.Error).Code, pq.ErrorCode("23505")) // is dublicate
			} else {
				require.NotEmpty(t, resp.ID.String)
				require.Empty(t, cmp.Diff(tt.input, resp, cmpopts.IgnoreFields(model.User{}, "ID")))
			}

			mustTruncateAll()
		})
	}
	mustTruncateAll()
}
