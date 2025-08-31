package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobAd struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AdTitle        string             `bson:"ad_title" json:"ad_title"`
	JobDescription string             `bson:"job_description" json:"job_description"`
	Qualification  string             `bson:"qualification" json:"qualification"`
	JobType        string             `bson:"job_type" json:"job_type"`
	JobId          primitive.ObjectID `bson:"job_id" json:"job_id"`
}

type JobAds []*JobAd

func (j *JobAds) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(j)
}

func (j *JobAd) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(j)
}

func (j *JobAd) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(j)
}
