package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	CitizenUCN     string             `bson:"citizen_ucn"`
	JobId          primitive.ObjectID `bson:"job_id"`
	StartDate      time.Time          `bson:"start_date"`
	EndDate        time.Time          `bson:"end_date"`
	EmployerReview string             `bson:"employer_review"`
}
