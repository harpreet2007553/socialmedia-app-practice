package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Follow struct {
	Followed primitive.ObjectID `json:"following" bson:"following"`
	FollowedBy primitive.ObjectID `json:"followed_by" bson:"followed_by"`
}