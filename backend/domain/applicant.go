package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Applicant struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	JobAdId primitive.ObjectID `bson:"job_ad_id" json:"job_ad_id"`
	CVId    primitive.ObjectID `bson:"cv_id,omitempty" json:"cv_id,omitempty"`
}

type Applicants []*Applicant

func (a *Applicants) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Applicant) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Applicant) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
