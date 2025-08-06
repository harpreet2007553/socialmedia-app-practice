package controllers

import (
	"backend-in-go/db"
	"backend-in-go/models"
	"backend-in-go/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Requested_User struct {
	FullName     string    `json:"fullName" bson:"fullName,omitempty"`
	UserName     string    `json:"UserName" bson:"UserName,omitempty"`
	Email        string    `json:"email" bson:"email,omitempty"`
	Avatar       string    `json:"avatar" bson:"avatar,omitempty"`
	Password     string    `json:"password" bson:"password,omitempty"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user Requested_User
	
	body, _ := io.ReadAll(r.Body)
    err := json.Unmarshal(body, &user)
    if err != nil {
        http.Error(w, "Invalid request payload!!", http.StatusBadRequest)
        return
    }

	
	// decoder := json.NewDecoder(r.Body)
    // decoder.DisallowUnknownFields()
    // err = decoder.Decode(&user)

	// if err != nil {
	// 	http.Error(w, "Invalid request payload", http.StatusBadRequest)
	// 	return
	// }
	defer r.Body.Close()
	checkUser := db.Collection_users.FindOne(context.TODO() , bson.M{
		"$or": []bson.M{
			{"username": user.UserName},
			{"email": user.Email},
		},
	})
	if checkUser.Err() == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	full_User := models.User{
		FullName:     user.FullName,
		UserName:     user.UserName,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Password:     user.Password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		RefreshToken: "<token>",
	}

	result, err := db.Collection_users.InsertOne(context.Background(), full_User)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	token_user := utils.JWTUser{
		Id: result.InsertedID.(string),
		UserName: user.UserName,
		Email: full_User.Email,
	}
	tokens,err := utils.GenerateJWT(token_user)
	if err != nil{
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}
	_ , err = db.Collection_users.UpdateOne(context.TODO(),bson.M{"_id": result.InsertedID}, bson.M{ "$set": bson.M{"RefreshToken": tokens.RefreshToken}} )

	if err != nil {
		http.Error(w, "failed to load Refresh Token", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message":    "User registered successfully!",
		"insertedId": result.InsertedID,
		"tokens":     tokens.RefreshToken,
	})
}