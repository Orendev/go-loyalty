package app

import (
	"encoding/json"
	"github.com/Orendev/go-loyalty/internal/auth"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

// PostWithdraw Запрос на списание средств
func (a *App) PostWithdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req := models.WithdrawRequest{}
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	now := time.Now()

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	number, err := strconv.Atoi(req.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account, err := a.repo.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if account.Current < req.Sum {
		http.Error(w, "", http.StatusPaymentRequired)
		return
	}

	transact := models.Transact{
		ID:          uuid.New().String(),
		Amount:      req.Sum,
		OrderNumber: number,
		AccountID:   account.ID,
		Debit:       false,
		ProcessedAt: now.Format(time.RFC3339),
	}

	if err = transact.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = a.repo.AddWithdraw(r.Context(), transact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	current := account.Current - req.Sum
	err = a.repo.UpdateAccountCurrent(r.Context(), account.ID, current)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// GetWithdraw Получение информации о выводе средств
func (a *App) GetWithdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limit := 100
	w.Header().Set("Content-Type", "application/json")
	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	account, err := a.repo.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	withdrawResponse := make([]models.WithdrawResponse, 0, limit)

	withdraws, err := a.repo.GetWithdrawByAccountID(r.Context(), account.ID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, withdraw := range withdraws {
		// заполняем модель ответа
		withdrawResponse = append(withdrawResponse, models.WithdrawResponse{
			Order:       strconv.Itoa(withdraw.OrderNumber),
			Sum:         withdraw.Amount,
			ProcessedAt: withdraw.ProcessedAt,
		})
	}

	// заполняем модель ответа
	enc, err := json.Marshal(withdrawResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(withdrawResponse) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, err = w.Write(enc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
