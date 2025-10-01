package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CitizenUCN     string             `bson:"citizen_ucn" json:"owner_ucn"`
	JobId          primitive.ObjectID `bson:"job_id" json:"job_id"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	EndDate        time.Time          `bson:"end_date,omitempty" json:"end_date,omitempty"`
	EmployerReview string             `bson:"employer_review,omitempty" json:"employer_review,omitempty"`
}
