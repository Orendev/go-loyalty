package app

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/Orendev/go-loyalty/internal/auth"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/google/uuid"
	"net/http"
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
	user := models.User{
		Login:    loginReq.Login,
		Password: hex.EncodeToString(hash[:]),
		Id:       uuid.New().String(),
	}
	err := a.repo.AddNewUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//аутентификация пользователя
	ctx, err := auth.NewSigner(r.Context())
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	r, err = auth.ContextToHTTP(w, r.WithContext(ctx))

	w.WriteHeader(http.StatusCreated)
}