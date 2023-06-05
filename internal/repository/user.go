package repository

import (
	"context"
	"fmt"
	"github.com/Orendev/go-loyalty/internal/models"
)

func (r *Repository) Login(ctx context.Context, login, password string) (u models.User, err error) {
	row := r.db.QueryRowContext(ctx, `select id, login from users where login = $1 AND password = $2`, login, password)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	err = row.Scan(&u.Id, &u.Login)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}

func (r *Repository) AddNewUser(ctx context.Context, u models.User) (err error) {
	_, err = r.db.ExecContext(ctx, `insert into users (id, login, password) values ($1, $2, $3)`, u.Id, u.Login, u.Password)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}
