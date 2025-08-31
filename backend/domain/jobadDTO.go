package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobAdDTO struct {
	JobAdID        primitive.ObjectID `bson:"_id" json:"_id"`
	AdTitle        string             `bson:"ad_title" json:"ad_title"`
	JobDescription string             `bson:"job_description" json:"job_description"`
	Qualification  string             `bson:"qualification" json:"qualification"`
	JobType        string             `bson:"job_type" json:"job_type"`
	CompanyName    string             `bson:"company_name" json:"company_name"`
	CompanyId      string             `bson:"company_id" json:"company_id"`
}
