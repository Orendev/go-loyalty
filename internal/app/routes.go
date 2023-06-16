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

		r.Route("/api/user", func(r chi.Router) {
			r.Post("/orders", a.PostOrders)
			r.Get("/orders", a.GetOrders)

			r.Post("/balance/withdraw", a.PostWithdraw)
			r.Get("/balance", a.GetBalance)

			r.Get("/withdrawals", a.GetWithdraw)
		})

	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Route("/api/user", func(r chi.Router) {
			r.Post("/register", a.Signup)
			r.Post("/login", a.Login)
		})
	})

	return r
}
