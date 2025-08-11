package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Content string             `json:"content" bson:"content,omitempty"`
	Owner   primitive.ObjectID `json:"owner" bson:"owner,omitempty"`
	PostID  primitive.ObjectID `json:"post_id" bson:"post_id,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}