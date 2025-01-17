package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewOfCompany struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Description string             `bson:"description"`
	Rating      int                `bson:"rating"`
	EmployeeId  primitive.ObjectID `bson:"employee_id"`
	EmployerId  primitive.ObjectID `bson:"employer_id"`
}

type ReviewsOfCompany []*ReviewOfCompany

func (rc *ReviewsOfCompany) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(rc)
}

func (rc *ReviewOfCompany) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(rc)
}

func (rc *ReviewOfCompany) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(rc)
}
