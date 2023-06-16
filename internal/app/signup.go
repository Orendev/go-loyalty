package app

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Orendev/go-loyalty/internal/auth"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/Orendev/go-loyalty/internal/repository"
	"github.com/google/uuid"
)

func (a *App) Signup(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var loginReq models.LoginRequest
	dec := json.NewDecoder(r.Body)
	// читаем тело запроса и декодируем
	if err := dec.Decode(&loginReq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := loginReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash := md5.Sum([]byte(loginReq.Password))
	userID := uuid.New().String()
	user := models.User{
		ID:       userID,
		Login:    loginReq.Login,
		Password: hex.EncodeToString(hash[:]),
	}
	err := a.repo.AddUser(r.Context(), user)
	if err != nil && errors.Is(err, repository.ErrorDuplicate) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	//аутентификация пользователя
	ctx, err := auth.NewSigner(r.Context(), userID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	_, err = auth.ContextToHTTP(w, r.WithContext(ctx))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
