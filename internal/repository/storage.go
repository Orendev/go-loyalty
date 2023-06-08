package repository

import (
	"context"
	"github.com/Orendev/go-loyalty/internal/models"
)

type Storage interface {
	Login(ctx context.Context, login, password string) (u models.User, err error)
	AddUser(ctx context.Context, u models.User) (err error)
	AddOrder(ctx context.Context, u models.Order) (err error)
	GetOrderByUserID(ctx context.Context, userID string, limit int) ([]models.Order, error)
	Ping(ctx context.Context) error
	Close() error
}
