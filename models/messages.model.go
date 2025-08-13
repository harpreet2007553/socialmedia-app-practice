package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Type   string             `bson:"type,omitempty"`
	Sender primitive.ObjectID `bson:"sender,omitempty"`
	Text   string             `bson:"text,omitempty"`
	Reciever primitive.ObjectID `bson:"reciever,omitempty"`
}