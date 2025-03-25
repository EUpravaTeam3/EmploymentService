package domain

import "time"

type User struct {
	UCN         string    `bson:"ucn"`
	Name        string    `bson:"name"`
	LastName    string    `bson:"lastName"`
	DateOfBirth time.Time `bson:"date_of_birth"`
	Address     string    `bson:"address"`
	Email       string    `bson:"email"`
	Password    string    `bson:"password"`
	Role        string    `bson:"role"`
}
