package usecases

import (
	"errors"
	"time"

	"github.com/yeahmerey/go-auth-service/internal/db"
)

type Gym = db.Gym

func GetGyms() []Gym {
	rows, _ := db.DB.Query("SELECT id, name, address, capacity, clients FROM gyms")
	defer rows.Close()
	var gyms []Gym
	for rows.Next() {
		var g Gym
		rows.Scan(&g.ID, &g.Name, &g.Address, &g.Capacity, &g.Clients)
		gyms = append(gyms, g)
	}
	return gyms
}
func JoinGym(userID, gymID int) error {
	// Проверяем, существует ли тренажерный зал
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM gyms WHERE id = $1)", gymID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("gym not found")
	}
	
	// Проверяем, не достигнута ли максимальная емкость зала
	var capacity, clients int
	err = db.DB.QueryRow("SELECT capacity, clients FROM gyms WHERE id = $1", gymID).Scan(&capacity, &clients)
	if err != nil {
		return err
	}
	
	if clients >= capacity {
		return errors.New("gym is at full capacity")
	}
	
	// Начинаем транзакцию
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Добавляем запись в gym_members
	_, err = tx.Exec(
		"INSERT INTO gym_members (user_id, gym_id, joined_at) VALUES ($1, $2, $3) ON CONFLICT (user_id, gym_id) DO NOTHING",
		userID, gymID, time.Now(),
	)
	if err != nil {
		return err
	}
	
	// Обновляем количество клиентов в зале
	_, err = tx.Exec("UPDATE gyms SET clients = clients + 1 WHERE id = $1", gymID)
	if err != nil {
		return err
	}
	
	// Фиксируем транзакцию
	return tx.Commit()
}
func LeaveGym(userID, gymID int) error {
	// Начинаем транзакцию
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Проверяем, является ли пользователь членом этого зала
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM gym_members WHERE user_id = $1 AND gym_id = $2)", userID, gymID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user is not a member of this gym")
	}
	
	// Удаляем запись из gym_members
	_, err = tx.Exec("DELETE FROM gym_members WHERE user_id = $1 AND gym_id = $2", userID, gymID)
	if err != nil {
		return err
	}
	
	// Обновляем количество клиентов в зале
	_, err = tx.Exec("UPDATE gyms SET clients = clients - 1 WHERE id = $1", gymID)
	if err != nil {
		return err
	}
	
	// Фиксируем транзакцию
	return tx.Commit()
}
func GetUserGyms(userID int) ([]db.Gym, error) {
	rows, err := db.DB.Query(`
		SELECT g.id, g.name, g.address, g.capacity, g.clients, gm.joined_at
		FROM gyms g
		JOIN gym_members gm ON g.id = gm.gym_id
		WHERE gm.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var gyms []db.Gym
	for rows.Next() {
		var g db.Gym
		var joinedAt time.Time
		err := rows.Scan(&g.ID, &g.Name, &g.Address, &g.Capacity, &g.Clients, &joinedAt)
		if err != nil {
			return nil, err
		}
		gyms = append(gyms, g)
	}
	
	return gyms, nil
}

func GetGymMembers(gymID int) ([]db.User, error) {
	rows, err := db.DB.Query(`
		SELECT u.id, u.username, u.email, gm.joined_at
		FROM users u
		JOIN gym_members gm ON u.id = gm.user_id
		WHERE gm.gym_id = $1
	`, gymID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []db.User
	for rows.Next() {
		var u db.User
		var joinedAt time.Time
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &joinedAt)
		if err != nil {
			return nil, err
		}
		// Не включаем пароль в результат
		u.Password = ""
		users = append(users, u)
	}
	
	return users, nil
}