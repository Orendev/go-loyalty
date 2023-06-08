package app

import (
	"encoding/json"
	"errors"
	"github.com/Orendev/go-loyalty/internal/auth"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/Orendev/go-loyalty/internal/repository"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"
	"time"
)

// PostOrders загрузка заказ
func (a *App) PostOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.GetAuthIdentifier(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	now := time.Now()
	number, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order := models.Order{
		ID:         uuid.New().String(),
		Number:     number,
		UserID:     userID,
		UploadedAt: now.Format(time.RFC3339),
	}

	if err = order.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = a.repo.AddOrder(r.Context(), order)
	if err != nil {
		if errors.Is(err, repository.ErrorDuplicate) {
			http.Error(w, err.Error(), http.StatusOK)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

// GetOrders получения списка загруженных заказов
func (a *App) GetOrders(w http.ResponseWriter, r *http.Request) {
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

	orderResponse := make([]models.OrderResponse, 0, limit)

	orders, err := a.repo.GetOrderByUserID(r.Context(), userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	for _, order := range orders {
		// заполняем модель ответа
		orderResponse = append(orderResponse, models.OrderResponse{
			Number:     strconv.Itoa(order.Number),
			Status:     order.Status,
			UploadedAt: order.UploadedAt,
		})
	}

	// заполняем модель ответа
	enc, err := json.Marshal(orderResponse)
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