package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Login(ctx context.Context, login, password string) (u models.User, err error) {
	row := r.db.QueryRowContext(ctx, `select id, login from users where login = $1 AND password = $2`, login, password)

	err = row.Scan(&u.ID, &u.Login)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}

func (r *Repository) AddNewUser(ctx context.Context, u models.User) (err error) {
	_, err = r.db.ExecContext(ctx, `insert into users (id, login, password) values ($1, $2, $3)`, u.ID, u.Login, u.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			err = ErrorDuplicate
		}
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}
	return
}
