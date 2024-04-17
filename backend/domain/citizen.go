package domain

import "time"

type Citizen struct {
	UCN         string    `json:"ucn"`
	Name        string    `json:"name"`
	LastName    string    `json:"lastName"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Address     string    `json:"address"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
}
