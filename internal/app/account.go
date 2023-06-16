package app

import (
	"encoding/json"
	"net/http"

	"github.com/Orendev/go-loyalty/internal/auth"
	"github.com/Orendev/go-loyalty/internal/models"
)

// GetBalance получение текущего баланса пользователя
func (a *App) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	account, err := a.repo.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// заполняем модель ответа
	accountResponse := models.AccountResponse{
		Current:   account.Current,
		Withdrawn: account.Withdrawn,
	}

	// заполняем модель ответа
	enc, err := json.Marshal(accountResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(enc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
