package main

import (
	"log"
	"net/http"

	"github.com/yeahmerey/go-auth-service/internal/db"
	"github.com/yeahmerey/go-auth-service/internal/handlers"
	"github.com/yeahmerey/go-auth-service/internal/middleware"
)

func main() {
	db.InitDB()

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/gyms", middleware.AuthMiddleware(handlers.GetGyms))

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}