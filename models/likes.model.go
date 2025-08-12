package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Like struct {
	ID primitive.ObjectID `json:"_id" bson:"_id, omitempty"`
	PostID primitive.ObjectID `json:"post_id" bson:"post_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"post_id"`
}