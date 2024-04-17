package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	Id          primitive.ObjectID `bson:"_id"`
	EmployerId  primitive.ObjectID `bson:"employer_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}
