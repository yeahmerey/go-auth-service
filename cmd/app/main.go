package main

import (
	"log"
	"net/http"
	"time"

	"github.com/yeahmerey/go-auth-service/internal/db"
	"github.com/yeahmerey/go-auth-service/internal/handlers"
	"github.com/yeahmerey/go-auth-service/internal/middleware"
	"github.com/yeahmerey/go-auth-service/internal/services"
)

func main() {
	db.InitDB()

    //clean of blacklist
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				services.CleanupBlacklist()
			}
		}
	}()

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", middleware.AuthMiddleware(handlers.Logout))
	
	http.HandleFunc("/gyms", middleware.AuthMiddleware(handlers.GetGyms))
	http.HandleFunc("/gyms/join", middleware.AuthMiddleware(handlers.JoinGym))
	http.HandleFunc("/gyms/leave", middleware.AuthMiddleware(handlers.LeaveGym))
	http.HandleFunc("/gyms/my", middleware.AuthMiddleware(handlers.GetUserGyms))
	http.HandleFunc("/gyms/members", middleware.AuthMiddleware(handlers.GetGymMembers))

	// runserver
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
