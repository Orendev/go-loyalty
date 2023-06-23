package middlewares

import (
	"net/http"

	"github.com/Orendev/go-loyalty/internal/auth"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter, http.Request как тот,
		// который будем передавать следующей функции
		ow := w
		or := r
		ctx, err := auth.HTTPToContext(or)
		if err != nil {
			http.Error(ow, "Unauthorized", http.StatusUnauthorized)
			return
		}

		or = or.WithContext(ctx)

		next.ServeHTTP(ow, or.WithContext(ctx))
	})
}
