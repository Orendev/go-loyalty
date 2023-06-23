package repository

import (
	"context"
	"fmt"

	"github.com/Orendev/go-loyalty/internal/logger"
	"github.com/Orendev/go-loyalty/internal/models"
	"go.uber.org/zap"
)

func (r *Repository) AddTransact(ctx context.Context, t models.Transact) error {
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
		`INSERT INTO transacts (id, amount, debit, order_number, account_id, processed_at) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, t.ID, t.Amount, t.Debit, t.OrderNumber, t.AccountID, t.ProcessedAt)

	if err != nil {
		return err
	}

	return tx.Commit()

}

func (r *Repository) GetWithdrawByAccountID(ctx context.Context, accountID string, limit int) ([]models.Transact, error) {

	transacts := make([]models.Transact, 0, limit)

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, amount, debit, order_number, account_id, processed_at
				FROM transacts
				WHERE account_id = $1 and debit=false
				ORDER BY processed_at
				`, accountID)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Log.Error("error", zap.Error(err))
		}
	}()

	// пробегаем по всем записям
	for rows.Next() {
		var t models.Transact
		err = rows.Scan(&t.ID, &t.Amount, &t.Debit, &t.OrderNumber, &t.AccountID, &t.ProcessedAt)
		if err != nil {
			err = fmt.Errorf("failed to query data: %w", err)
			return nil, err
		}

		transacts = append(transacts, t)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return transacts, nil
}
