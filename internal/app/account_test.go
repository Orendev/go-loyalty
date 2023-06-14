package app_test

import (
	"context"
	"fmt"
	"github.com/Orendev/go-loyalty/internal/app"
	"github.com/Orendev/go-loyalty/internal/auth"
	"github.com/Orendev/go-loyalty/internal/middlewares"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/Orendev/go-loyalty/internal/repository/mock"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApp_GetBalance(t *testing.T) {

	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mock.NewMockStorage(ctrl)

	now := time.Now()

	a, err := app.NewApp(context.Background(), s, make(chan models.Accrual, 10))
	if err != nil {
		require.NoError(t, err)
	}

	uri := "/api/user/balance"
	r := chi.NewRouter()
	r.Use(middlewares.Auth)
	r.Get(uri, a.GetBalance)

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		contentType  string
		expectedCode int
		expectedBody string
	}

	type args struct {
		userID       string
		accountModel *models.Account
	}

	tests := []struct {
		name   string // добавляем название тестов
		method string
		args   args
		want   want
	}{
		{
			name:   "method_get_unauthorized",
			method: http.MethodGet,
			want: want{
				expectedCode: http.StatusUnauthorized,
			},
		},
		{
			name:   "method_get_success",
			method: http.MethodGet,
			args: args{
				userID: uuid.New().String(),
				accountModel: &models.Account{
					ID:        uuid.New().String(),
					Current:   50050,
					Withdrawn: 500,
					CreatedAt: now.Format(time.RFC3339),
					UpdatedAt: now.Format(time.RFC3339),
				},
			},
			want: want{
				contentType:  "application/json",
				expectedCode: http.StatusOK,
				expectedBody: `{"current":500.5,"withdrawn":5}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.args.accountModel != nil {
				s.EXPECT().
					GetAccountByUserID(gomock.Any(), gomock.Any()).
					Return(tt.args.accountModel, nil)
			}

			var bodyReader io.Reader
			req, err := http.NewRequest(tt.method, srv.URL+uri, bodyReader)
			require.NoError(t, err)

			if tt.args.userID != "" {
				ctx, err := auth.NewSigner(context.Background(), tt.args.userID)
				require.NoError(t, err)

				token, ok := ctx.Value(auth.JwtContextKey).(string)
				if !ok {
					require.NoError(t, auth.ErrorTokenContextMissing)
				}

				req.Header.Set(auth.HeaderAuthorizationKey, fmt.Sprintf("Bearer %s", token))
			}

			resp, err := srv.Client().Do(req)
			require.NoError(t, err)

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}()

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.want.expectedCode, resp.StatusCode, "code didn't match expected")

			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "header content type didn't match expected")
			}

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tt.want.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					require.NoError(t, err)
				}
				assert.Regexp(t, tt.want.expectedBody, string(body))
			}
		})
	}
}
