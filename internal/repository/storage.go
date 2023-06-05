package repository

import (
	"context"
	"github.com/Orendev/go-loyalty/internal/models"
)

type Storage interface {
	Login(ctx context.Context, login, password string) (u models.User, err error)
	AddNewUser(ctx context.Context, u models.User) (err error)
	Ping(ctx context.Context) error
	Close() error
}
