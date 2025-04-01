package db

import "time"

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type Gym struct {
	ID       int
	Name     string
	Address  string
	Capacity int
	Clients  int
}

type GymMember struct {
	ID       int
	UserID   int
	GymID    int
	JoinedAt time.Time
}
