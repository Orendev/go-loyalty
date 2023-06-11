package repository

import (
	"context"
	"github.com/Orendev/go-loyalty/internal/models"
)

func (r *Repository) AddWithdraw(ctx context.Context, t models.Transact) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO transacts (id, amount, debit, order_number, account_id, processed_at) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, t.ID, t.Amount, t.Debit, t.OrderNumber, t.AccountID, t.ProcessedAt)

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
