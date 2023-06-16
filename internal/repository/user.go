package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Orendev/go-loyalty/internal/logger"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *Repository) Login(ctx context.Context, login, password string) (models.User, error) {
	u := models.User{}
	row := r.db.QueryRowContext(ctx, `select id, login from users where login = $1 AND password = $2`, login, password)

	err := row.Scan(&u.ID, &u.Login)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return u, err
	}
	return u, nil
}

func (r *Repository) AddUser(ctx context.Context, u models.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		// если ошибка, то откатываем изменения
		err = tx.Rollback()
		if err != nil {
			logger.Log.Error("error", zap.Error(err))
		}
	}()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, u.ID, u.Login, u.Password)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	now := time.Now()
	timestamp := now.Format(time.RFC3339)

	account := models.Account{
		ID:        uuid.New().String(),
		UserID:    u.ID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	return r.AddAccount(ctx, account)
}
