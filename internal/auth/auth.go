package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type contextKey string

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
// Её встраивают в структуру утверждений, определённую пользователем.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const (
	TokenExp                          = time.Hour * 3
	SecretKey                         = "supersecretkey"
	JwtContextKey          contextKey = "JWTToken"
	JwtUserIDContextKey    contextKey = "JWTUserID"
	HeaderAuthorizationKey string     = "Authorization"
)

var (
	// ErrorTokenContextMissing токен не был передан
	ErrorTokenContextMissing = errors.New("token up for parsing was not passed through the context")

	// ErrorTokenInvalid означает, что токен не удалось проверить.
	ErrorTokenInvalid = errors.New("JWT was invalid")

	// ErrorUnexpectedSigningMethod означает, что токен был подписан с использованием неожиданного метода подписи.
	ErrorUnexpectedSigningMethod = errors.New("unexpected signing method")

	// ErrorTokenMalformed токен не был отформатирован как JWT.
	ErrorTokenMalformed = errors.New("JWT is malformed")

	// ErrorTokenExpired заголовок срока действия токена прошел.
	ErrorTokenExpired = errors.New("JWT is expired")

	// ErrorTokenNotActive Токен еще не действителен
	ErrorTokenNotActive = errors.New("token is not valid yet")
)

const (
	bearer       string = "bearer"
	bearerFormat string = "Bearer %s"
)

func extractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], bearer) {
		return "", false
	}

	return authHeaderParts[1], true
}

func newParse(ctx context.Context) (context.Context, error) {
	claims := &Claims{}
	tokenString, ok := ctx.Value(JwtContextKey).(string)
	if !ok {
		return nil, ErrorTokenContextMissing
	}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrorUnexpectedSigningMethod
			}
			return []byte(SecretKey), nil
		})
	//*jwt.ValidationError
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			// Token is malformed
			return nil, ErrorTokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			// Token is expired
			return context.WithValue(ctx, JwtUserIDContextKey, claims.UserID), ErrorTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			// Token is not active yet
			return nil, ErrorTokenNotActive
		}

		return nil, err
	}

	if !token.Valid {
		return nil, ErrorTokenInvalid
	}

	return context.WithValue(ctx, JwtUserIDContextKey, claims.UserID), nil
}

func generateAuthHeaderFromToken(token string) string {
	return fmt.Sprintf(bearerFormat, token)
}

func GetAuthIdentifier(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(JwtUserIDContextKey).(string)
	if !ok {
		return "", ErrorTokenContextMissing
	}
	return userID, nil
}

func HTTPToContext(r *http.Request) (context.Context, error) {
	token, ok := extractTokenFromAuthHeader(r.Header.Get(HeaderAuthorizationKey))
	if !ok {
		return nil, ErrorTokenContextMissing
	}

	return newParse(context.WithValue(r.Context(), JwtContextKey, token))
}

func ContextToHTTP(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	ctx := r.Context()
	tokenString, ok := ctx.Value(JwtContextKey).(string)

	if !ok {
		return nil, ErrorTokenContextMissing
	}

	w.Header().Add(HeaderAuthorizationKey, generateAuthHeaderFromToken(tokenString))

	return r, nil
}

// NewSigner создаёт JWT, указывая идентификатор ключа,
func NewSigner(ctx context.Context) (context.Context, error) {

	userID, err := GetAuthIdentifier(ctx)
	if err != nil {
		userID = uuid.New().String()
	}
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, JwtContextKey, tokenString), nil
}
