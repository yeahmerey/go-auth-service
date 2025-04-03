package usecases

import (
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
