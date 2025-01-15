package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobAd struct {
	Id             primitive.ObjectID `bson:"job_ad_id,omitempty"`
	AdTitle        string             `bson:"ad_title"`
	JobDescription string             `bson:"job_description"`
	Qualification  string             `bson:"qualification"`
	JobType        string             `bson:"job_type"`
	JobId          primitive.ObjectID `bson:"job_id"`
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
