package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Orendev/go-loyalty/internal/logger"
	"github.com/Orendev/go-loyalty/internal/models"
	"go.uber.org/zap"
)

func (r *Repository) AddAccount(ctx context.Context, a models.Account) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			logger.Log.Error("error", zap.Error(err))
		}
	}()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO accounts (id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, a.ID, a.UserID, a.CreatedAt, a.UpdatedAt)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetAccountByUserID(ctx context.Context, userID string) (*models.Account, error) {

	row := r.db.QueryRowContext(ctx,
		`SELECT a.id, a.current, a.user_id, a.created_at, a.updated_at, sum(t.amount) as withdrawn
				FROM accounts a
					LEFT JOIN transacts t ON a.id = t.account_id AND t.debit=false
				WHERE user_id = $1
				GROUP BY a.id`, userID)

	account := models.Account{}

	var withdrawn sql.NullInt64

	err := row.Scan(&account.ID, &account.Current, &account.UserID, &account.CreatedAt, &account.UpdatedAt, &withdrawn)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return nil, err
	}

	if withdrawn.Valid {
		account.Withdrawn = int(withdrawn.Int64)
	}

	return &account, nil
}

func (r *Repository) UpdateAccountCurrent(ctx context.Context, id string) error {
	current, err := r.GetCurrent(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err = tx.Rollback()
		if err != nil {
			logger.Log.Error("error", zap.Error(err))
		}
	}()

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE accounts SET current = $1, updated_at = $2 WHERE id = $3`)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	now := time.Now()
	_, err = stmt.ExecContext(ctx, current, now.Format(time.RFC3339), id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetCurrent(ctx context.Context, accountID string) (int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT sum(t.amount) AS amount, debit
				FROM transacts t
				WHERE t.account_id = $1
				GROUP BY t.debit
				ORDER BY t.debit DESC
				`, accountID)
	if err != nil {
		return 0, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Log.Error("error", zap.Error(err))
		}
	}()

	var amount int
	var debit bool
	var current int

	current = 0

	// пробегаем по всем записям
	for rows.Next() {
		err = rows.Scan(&amount, &debit)
		if err != nil {
			err = fmt.Errorf("failed to query data: %w", err)
			return 0, err
		}

		if debit {
			current += amount
		} else {
			current -= amount
		}
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return current, nil
}
