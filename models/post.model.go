package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
	Attachment string `json:"attachment" bson:"attachment,omitempty"`
	Title string `json:"title" bson:"title,omitempty"`
	Owner primitive.ObjectID `json:"owner" bson:"owner,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}
