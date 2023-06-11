package repository

import (
	"context"
	"fmt"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/google/uuid"
	"time"
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

func (r *Repository) AddUser(ctx context.Context, u models.User) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`)
	if err != nil {
		return
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, u.ID, u.Login, u.Password)
	if err != nil {
		// если ошибка, то откатываем изменения
		errRollback := tx.Rollback()
		if errRollback != nil {
			err = errRollback
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	now := time.Now()
	timestamp := now.Format(time.RFC3339)

	account := models.Account{
		ID:        uuid.New().String(),
		UserID:    u.ID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	err = r.AddAccount(ctx, account)

	return
}
