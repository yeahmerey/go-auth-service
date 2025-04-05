
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
	// Инициализация базы данных
	db.InitDB()

	// Запускаем очистку черного списка токенов каждый час
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

	// Регистрация обработчиков запросов
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", middleware.AuthMiddleware(handlers.Logout)) // Добавлен middleware для проверки авторизации
	http.HandleFunc("/gyms", middleware.AuthMiddleware(handlers.GetGyms))

	// Запуск сервера
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}