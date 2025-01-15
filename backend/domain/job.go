package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	PoistionName     string             `bson:"position_name"`
	Pay              int                `bson:"pay"`
	EmployerId       primitive.ObjectID `bson:"employer_id"`
	NumOfEmployees   int                `bson:"num_of_employees"`
	EmployeeCapacity int                `bson:"employee_capacity"`
}

type Jobs []*Job

func (j *Jobs) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(j)
}

func (j *Job) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(j)
}

func (j *Job) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(j)
}
