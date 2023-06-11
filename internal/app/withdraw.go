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
