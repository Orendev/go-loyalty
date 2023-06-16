package repository

import (
	"context"

	"github.com/Orendev/go-loyalty/internal/models"
)

type Storage interface {
	Login(ctx context.Context, login, password string) (u models.User, err error)
	AddUser(ctx context.Context, u models.User) (err error)
	AddAccount(ctx context.Context, a models.Account) (err error)
	AddOrder(ctx context.Context, u models.Order) (err error)
	AddTransact(ctx context.Context, t models.Transact) (err error)
	UpdateAccountCurrent(ctx context.Context, id string) (err error)
	GetCurrent(ctx context.Context, accountID string) (int, error)
	GetOrderByUserID(ctx context.Context, userID string, limit int) ([]models.Order, error)
	GetOrderByNumber(ctx context.Context, number int, userID string) (models.Order, error)
	UpdateStatusOrder(ctx context.Context, orderNumber int, status string) error
	GetAccountByUserID(ctx context.Context, userID string) (*models.Account, error)
	GetWithdrawByAccountID(ctx context.Context, accountID string, limit int) ([]models.Transact, error)
	Ping(ctx context.Context) error
	Close() error
}
