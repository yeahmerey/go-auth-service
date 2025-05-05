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
	// get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// check token and get username
	claims, err := services.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// get user ID from database
	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// get gym ID from request
	var req GymMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// add user to gym
	if err := usecases.JoinGym(userID, req.GymID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully joined gym"})
}
func LeaveGym(w http.ResponseWriter, r *http.Request) {
	// get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// check token and get username
	claims, err := services.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	//get user id from database
	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// get gym ID from request
	var req GymMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// delete user from gym
	if err := usecases.LeaveGym(userID, req.GymID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully left gym"})
}
func GetUserGyms(w http.ResponseWriter, r *http.Request) {
	// get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// check token and get username
	claims, err := services.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// get user ID from database
	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// get list of gyms for the user
	gyms, err := usecases.GetUserGyms(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gyms)
}
func GetGymMembers(w http.ResponseWriter, r *http.Request) {
	// get id of the gym from the request
	gymIDStr := r.URL.Query().Get("gym_id")
	gymID, err := strconv.Atoi(gymIDStr)
	if err != nil {
		http.Error(w, "Invalid gym ID", http.StatusBadRequest)
		return
	}

	// get list of users in the gym
	users, err := usecases.GetGymMembers(gymID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
