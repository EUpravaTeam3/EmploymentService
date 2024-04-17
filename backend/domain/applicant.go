package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Applicant struct {
	Id      primitive.ObjectID `bson:"_id"`
	JobAdId primitive.ObjectID `bson:"job_ad_id"`
	CVId    primitive.ObjectID `bson:"cv_id"`
}
