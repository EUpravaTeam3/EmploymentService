package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CV struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CitizenUCN     string             `bson:"citizen_ucn" json:"citizen_ucn"`
	Name           string             `bson:"name" json:"name"`
	Email          string             `bson:"email" json:"email"`
	Description    string             `bson:"description" json:"description"`
	WorkExperience []string           `bson:"work_experience" json:"work_experience"`
	Education      []Diploma          `bson:"education" json:"education"`
}

type Cvs []*CV

func (c *Cvs) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *CV) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *CV) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}
