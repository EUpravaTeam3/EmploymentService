package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Job struct {
	Id               primitive.ObjectID `bson:"_id"`
	PoistionName     string             `bson:"position_name"`
	Pay              int                `bson:"pay"`
	EmployerId       primitive.ObjectID `bson:"employer_id"`
	NumOfEmployees   int                `bson:"num_of_employees"`
	EmployeeCapacity int                `bson:"employee_capacity"`
}
