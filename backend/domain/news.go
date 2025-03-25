package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type News struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	CompanyId   primitive.ObjectID `bson:"employer_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
}

type AllNews []*Job

func (j *News) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(j)
}

func (j *AllNews) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(j)
}

func (j *News) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(j)
}
