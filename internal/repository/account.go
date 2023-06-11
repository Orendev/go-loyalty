package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Orendev/go-loyalty/internal/models"
	"time"
)

func (r *Repository) AddAccount(ctx context.Context, a models.Account) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO accounts (id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, a.ID, a.UserID, a.CreatedAt, a.UpdatedAt)

	if err != nil {
		// если ошибка, то откатываем изменения
		errRollback := tx.Rollback()
		if errRollback != nil {
			err = errRollback
		}
		return
	}

	err = tx.Commit()
	return
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

func (r *Repository) UpdateAccountCurrent(ctx context.Context, id string, current int) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE accounts SET current = $1, updated_at = $2 WHERE id = $3`)
	if err != nil {
		return
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	now := time.Now()
	updatedAt := now.Format(time.RFC3339)
	_, err = stmt.ExecContext(ctx, current, updatedAt, id)

	if err != nil {
		// если ошибка, то откатываем изменения
		errRollback := tx.Rollback()
		if errRollback != nil {
			err = errRollback
		}
		return
	}

	err = tx.Commit()
	return
}
