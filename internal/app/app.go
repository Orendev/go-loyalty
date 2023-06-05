package app

import (
	"context"
	"github.com/Orendev/go-loyalty/internal/middlewares"
	"github.com/Orendev/go-loyalty/internal/repository"
	"github.com/go-chi/chi/v5"
)

type App struct {
	repo repository.Storage
}

func (a *App) Routes(r chi.Router) chi.Router {

	r.Use(middlewares.Logger)
	r.Use(middlewares.Gzip)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		//r.Post("/api/user/orders", a.Signup)
		//r.Get("/api/user/orders", a.Signup)
		//r.Get("/api/user/balance", a.Login)
		//r.Post("/api/user/balance/withdraw", a.Login)
		//r.Get("/api/user/withdrawals", a.Login)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", a.Signup)
		r.Post("/api/user/login", a.Login)
	})

	return r
}

func NewApp(_ context.Context, repo repository.Storage) (*App, error) {

	return &App{repo: repo}, nil
}
