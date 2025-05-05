package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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
	// get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}
	
	// delete if token is in the blacklist
	token := strings.TrimPrefix(authHeader, "Bearer ")
	
	// add token to the blacklist
	if err := usecases.Logout(token); err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
	
	// send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}
