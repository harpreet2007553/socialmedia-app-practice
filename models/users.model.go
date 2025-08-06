package models

import "time"

type User struct {

	FullName     string    `json:"fullName" bson:"fullName,omitempty"`
	UserName     string    `json:"UserName" bson:"UserName,omitempty"`
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
	UserName     string    `json:"UserName" bson:"UserName,omitempty"`
	Email        string    `json:"email" bson:"email,omitempty"`
}


