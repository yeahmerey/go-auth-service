package usecases

import (
	"errors"
	"github.com/yeahmerey/go-auth-service/internal/db"
	"github.com/yeahmerey/go-auth-service/internal/services"
)

func Register(username, email, password string) error {
	hash, err := services.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", username, email, hash)
	return err
}

func Login(username, password string) (map[string]string, error) {
	row := db.DB.QueryRow("SELECT password FROM users WHERE username=$1", username)
	var hashed string
	if err := row.Scan(&hashed); err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := services.CheckPassword(hashed, password); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return services.GenerateTokens(username)
}
func Logout(token string) error {
	return services.BlacklistToken(token)
}
