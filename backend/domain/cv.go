package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CV struct {
	Id             primitive.ObjectID `bson:"_id"`
	CitizenUCN     string             `bson:"citizen_ucn"`
	Description    string             `bson:"description"`
	WorkExperience []string           `bson:"work_experience"`
	Education      []Diploma          `bson:"education"`
}
