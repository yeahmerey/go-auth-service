package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/yeahmerey/go-auth-service/internal/usecases"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)
	if err := usecases.Register(creds.Username, creds.Email, creds.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)
	tokens, err := usecases.Login(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(tokens)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// В базовом варианте просто вернём 200 OK
	w.WriteHeader(http.StatusOK)
}