package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/yeahmerey/go-auth-service/internal/usecases"
)

func GetGyms(w http.ResponseWriter, r *http.Request) {
	gyms := usecases.GetGyms()
	json.NewEncoder(w).Encode(gyms)
}