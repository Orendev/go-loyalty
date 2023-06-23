package app

import (
	"github.com/Orendev/go-loyalty/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func (a *App) Routes(r chi.Router) chi.Router {

	r.Use(middlewares.Logger)
	r.Use(middlewares.Gzip)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Post("/api/user/orders", a.PostOrders)
		r.Get("/api/user/orders", a.GetOrders)
		r.Get("/api/user/balance", a.GetBalance)
		r.Post("/api/user/balance/withdraw", a.PostWithdraw)
		r.Get("/api/user/withdrawals", a.GetWithdraw)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", a.Signup)
		r.Post("/api/user/login", a.Login)
	})

	return r
}
