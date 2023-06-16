package app

import (
	"context"

	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/Orendev/go-loyalty/internal/repository"
)

type App struct {
	repo        repository.Storage
	accrualChan chan models.Accrual
}

func NewApp(_ context.Context, repo repository.Storage, accrualChan chan models.Accrual) *App {
	return &App{repo: repo, accrualChan: accrualChan}
}
