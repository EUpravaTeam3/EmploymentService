package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReviewOfCompany struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Description string             `bson:"description"`
	Rating      int                `bson:"rating"`
	EmployeeId  primitive.ObjectID `bson:"employee_id"`
	EmployerId  primitive.ObjectID `bson:"employer_id"`
}
