package client

import (
	"context"
	"time"

	"github.com/Orendev/go-loyalty/internal/logger"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *HTTPClient) worker(ctx context.Context) {

	logger.Log.Info("begin worker")

	for accrual := range h.accrualChan {

		logger.Log.Info("async go", zap.Int("order", accrual.Order))

		accrualResponse, err := h.GetAccrual(accrual.Order)
		if err != nil {
			logger.Log.Error("Error", zap.Any("GetAccrual", err))
			continue
		}

		logger.Log.Info("async go getAccrual", zap.Any("accrualResponse", accrualResponse))

		switch accrualResponse.Status {
		case models.StatusAccrualProcessed:
			accrual.Accrual = accrualResponse.Accrual
			err = h.addAccrual(ctx, accrual)
			if err != nil {
				logger.Log.Error("Error", zap.Any(models.StatusAccrualProcessed, err))
				continue
			}
		case models.StatusAccrualInvalid:
			err = h.repo.UpdateStatusOrder(ctx, accrual.Order, models.StatusOrderInvalid)
			if err != nil {
				logger.Log.Error("Error", zap.Any(models.StatusAccrualInvalid, err))
				continue
			}

		case models.StatusAccrualProcessing:
			err = h.repo.UpdateStatusOrder(ctx, accrual.Order, models.StatusOrderProcessing)
			if err != nil {
				logger.Log.Error("Error", zap.Any(models.StatusAccrualProcessing, err))
				continue
			}
			h.accrualChan <- models.Accrual{Order: accrual.Order, UserID: accrual.UserID}
		case models.StatusAccrualRegistered:
			h.accrualChan <- models.Accrual{Order: accrual.Order, UserID: accrual.UserID}
		default:
			h.accrualChan <- models.Accrual{Order: accrual.Order, UserID: accrual.UserID}
		}
		if accrualResponse.RetryAfterDuration != 0 {
			timer := time.NewTimer(accrualResponse.RetryAfterDuration * time.Second) // создаём таймер
			<-timer.C                                                                // ожидаем срабатывания таймера
		}

	}

	logger.Log.Info("end worker")
}

func (h *HTTPClient) addAccrual(ctx context.Context, accrual models.Accrual) error {
	account, err := h.repo.GetAccountByUserID(ctx, accrual.UserID)
	if err != nil {
		return err
	}

	now := time.Now()

	transact := models.Transact{
		ID:          uuid.New().String(),
		Amount:      accrual.Accrual,
		OrderNumber: accrual.Order,
		AccountID:   account.ID,
		Debit:       true,
		ProcessedAt: now.Format(time.RFC3339),
	}

	err = h.repo.AddTransact(ctx, transact)
	if err != nil {
		return err
	}

	err = h.repo.UpdateAccountCurrent(ctx, account.ID)
	if err != nil {
		return err
	}

	err = h.repo.UpdateStatusOrder(ctx, accrual.Order, models.StatusOrderProcessed)
	if err != nil {
		return err
	}

	return nil
}
