package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobAd struct {
	Id             primitive.ObjectID `bson:"job_ad_id,omitempty"`
	AdTitle        string             `bson:"ad_title"`
	JobDescription string             `bson:"job_description"`
	Qualification  string             `bson:"qualification"`
	JobType        string             `bson:"job_type"`
	JobId          primitive.ObjectID `bson:"job_id"`
}
