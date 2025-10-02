package domain

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CitizenUCN     string             `bson:"citizen_ucn" json:"citizen_ucn"`
	JobId          primitive.ObjectID `bson:"job_id" json:"job_id"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	EndDate        time.Time          `bson:"end_date,omitempty" json:"end_date,omitempty"`
	EmployerReview string             `bson:"employer_review,omitempty" json:"employer_review,omitempty"`
}

type Employees []*Employee

func (e *Employees) ToJSON(w io.Writer) error {
	en := json.NewEncoder(w)
	return en.Encode(e)
}

func (e *Employee) ToJSON(w io.Writer) error {
	en := json.NewEncoder(w)
	return en.Encode(e)
}

func (e *Employee) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(e)
}
