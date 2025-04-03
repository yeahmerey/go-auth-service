package db

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