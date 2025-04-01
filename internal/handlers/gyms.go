package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/yeahmerey/go-auth-service/internal/db"
	"github.com/yeahmerey/go-auth-service/internal/services"
	"github.com/yeahmerey/go-auth-service/internal/usecases"
)

type GymMemberRequest struct {
	GymID int `json:"gym_id"`
}

func GetGyms(w http.ResponseWriter, r *http.Request) {
	gyms := usecases.GetGyms()
	json.NewEncoder(w).Encode(gyms)
}

func JoinGym(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из заголовка
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	
	// Проверяем токен и получаем username
	claims, err := services.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	
	// Получаем ID пользователя из базы данных
	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	
	// Получаем ID тренажерного зала из запроса
	var req GymMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Добавляем пользователя в тренажерный зал
	if err := usecases.JoinGym(userID, req.GymID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully joined gym"})
}
func LeaveGym(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из заголовка
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	
	// Проверяем токен и получаем username
	claims, err := services.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	
	// Получаем ID пользователя из базы данных
	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	
	// Получаем ID тренажерного зала из запроса
	var req GymMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Удаляем пользователя из тренажерного зала
	if err := usecases.LeaveGym(userID, req.GymID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully left gym"})
}
func GetUserGyms(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из заголовка
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	
	// Проверяем токен и получаем username
	claims, err := services.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	
	// Получаем ID пользователя из базы данных
	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	
	// Получаем список тренажерных залов пользователя
	gyms, err := usecases.GetUserGyms(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gyms)
}
func GetGymMembers(w http.ResponseWriter, r *http.Request) {
	// Получаем ID тренажерного зала из URL параметров
	gymIDStr := r.URL.Query().Get("gym_id")
	gymID, err := strconv.Atoi(gymIDStr)
	if err != nil {
		http.Error(w, "Invalid gym ID", http.StatusBadRequest)
		return
	}
	
	// Получаем список пользователей тренажерного зала
	users, err := usecases.GetGymMembers(gymID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}