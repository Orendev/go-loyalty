package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Orendev/go-loyalty/internal/logger"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (r *Repository) AddOrder(ctx context.Context, o models.Order) error {

	_, err := r.db.ExecContext(ctx, `insert into orders (id, number, user_id, uploaded_at) values ($1, $2, $3, $4)`, o.ID, o.Number, o.UserID, o.UploadedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			err = ErrorDuplicate
		}
		err = fmt.Errorf("failed to exec data: %w", err)
		return err
	}
	return err
}

func (r *Repository) GetOrderByNumber(ctx context.Context, number int, userID string) (models.Order, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT o.id, o.number, o.status, o.user_id, o.uploaded_at, t.amount AS accrual
				FROM orders o
					LEFT JOIN transacts t ON o.number = t.order_number AND t.debit=true
				WHERE o.number = $1 AND user_id = $2`, number, userID)

	order := models.Order{}

	var accrual sql.NullInt64

	err := row.Scan(&order.ID, &order.Number, &order.Status, &order.UserID, &order.UploadedAt, &accrual)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return order, err
	}

	if accrual.Valid {
		order.Accrual = int(accrual.Int64)
	}

	return order, nil
}

func (r *Repository) GetOrderByUserID(ctx context.Context, userID string, limit int) ([]models.Order, error) {

	orders := make([]models.Order, 0, limit)

	rows, err := r.db.QueryContext(ctx,
		`SELECT o.id, o.number, o.status, o.user_id, o.uploaded_at, t.amount AS accrual
				FROM orders o
					LEFT JOIN transacts t ON o.number = t.order_number AND t.debit=true
				WHERE user_id = $1
				ORDER BY uploaded_at
				`, userID)
	if err != nil {
		return nil, err
	}

	var accrual sql.NullInt64

	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Log.Error("error", zap.Error(err))
		}
	}()

	// пробегаем по всем записям
	for rows.Next() {
		var o models.Order
		err = rows.Scan(&o.ID, &o.Number, &o.Status, &o.UserID, &o.UploadedAt, &accrual)
		if err != nil {
			err = fmt.Errorf("failed to query data: %w", err)
			return nil, err
		}

		if accrual.Valid {
			o.Accrual = int(accrual.Int64)
		}

		orders = append(orders, o)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) UpdateStatusOrder(ctx context.Context, orderNumber int, status string) error {
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
		`UPDATE orders SET status = $1 WHERE number = $2`)
	if err != nil {
		return err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			return
		}
	}()

	_, err = stmt.ExecContext(ctx, status, orderNumber)

	if err != nil {
		return err
	}

	return tx.Commit()
}
