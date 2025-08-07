package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FullName     string    `json:"fullName" bson:"fullName,omitempty"`
	UserName     string    `json:"userName" bson:"userName,omitempty"`
	Email        string    `json:"email" bson:"email,omitempty"`
	Avatar       string    `json:"avatar" bson:"avatar,omitempty"`
	Password     string    `json:"password" bson:"password,omitempty"`
	RefreshToken string    `json:"refreshToken" bson:"refreshToken,omitempty"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}
type UserResponse struct {
	Id           string    `json:"id" bson:"_id,omitempty"`
	FullName     string    `json:"fullName" bson:"fullName,omitempty"`
	UserName     string    `json:"userName" bson:"userName,omitempty"`
	Email        string    `json:"email" bson:"email,omitempty"`
}


