package controllers

import (
	// "backend-in-go/CyclicPackagesImport"
	"backend-in-go/db"
	"backend-in-go/middlewares"

	// "backend-in-go/middlewares"
	"backend-in-go/models"
	"backend-in-go/utils"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	// "log"
	"net/http"
	"time"

	// "backend-in-go/controllers/cookies"
	// "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Requested_User struct {
	FullName     string    `json:"fullName" bson:"fullName,omitempty"`
	UserName     string    `json:"userName" bson:"userName,omitempty"`
	Email        string    `json:"email" bson:"email,omitempty"`
	// Avatar       string    `json:"avatar" bson:"avatar,omitempty"`
	Password     string    `json:"password" bson:"password,omitempty"`
}
type Register_User_Cookie struct {
		RefreshToken string
		AccessToken string
}

// type contextKey struct{}


func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user Requested_User
	
	// body, _ := io.ReadAll(r.Body)
    // err := json.Unmarshal(body, &user)
    // if err != nil {
    //     http.Error(w, "Invalid request payload!!", http.StatusBadRequest)
    //     return
    // }

	
	decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    err := decoder.Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
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
		// Avatar:       user.Avatar,
		Password:     user.Password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result, err := db.Collection_users.InsertOne(context.Background(), full_User)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	InsertedId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	fmt.Println("User registered successfully with ID:", InsertedId.Hex())
	idStr := InsertedId.Hex()
	token_user := utils.JWTUser{
		Id: idStr,
		UserName: user.UserName,
		Email: full_User.Email,
	}
	
	tokens,err := utils.GenerateJWT(token_user)
	if err != nil{
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}
	_ , err = db.Collection_users.UpdateOne(context.TODO(),bson.M{"_id": InsertedId}, bson.M{ "$set": bson.M{"refreshToken": tokens.RefreshToken}} )

	if err != nil {
		http.Error(w, "failed to load Refresh Token", 500)
	}
    
	// TODO: "Passing Data Into Cookies"
    
    var cookie_data = Register_User_Cookie{
		RefreshToken:     tokens.RefreshToken,
		AccessToken:     tokens.AccessToken,
	}
	var buf bytes.Buffer;
    
	err = gob.NewEncoder(&buf).Encode(&cookie_data)
	if err != nil {
		http.Error(w, "Failed to encode cookie data", http.StatusInternalServerError)
		return
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())


	cookie := &http.Cookie{
		Name:     "user_JWT",
		Value:   encoded,
		Path: "/",
		Expires: time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Secure: false, // Set to true if using HTTPS

	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "User registered successfully!",
		"insertedId": result.InsertedID,
		"RefreshToken":     tokens.RefreshToken,
		"AccessToken":     tokens.AccessToken,
	})

}

type User_Login struct{
	Username string `json:"userName"`
	Password string `json:"password"`
	Email string `json:"email"`
}

var login_user User_Login
var user models.User

func LoginUser(w http.ResponseWriter, r *http.Request) {
	

	json.NewDecoder(r.Body).Decode(&login_user)

	if login_user.Username == "" && login_user.Email == "" {
		http.Error(w, "Username or Email is required", http.StatusBadRequest)
		return
	}
	

	if strings.TrimSpace(login_user.Email) == "" {
		// fmt.Printf("username type : %T",login_user.Username)
        err := db.Collection_users.FindOne(context.TODO(), bson.M{
				"userName": login_user.Username,
			}).Decode(&user)

		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			log.Fatal("User not found:", err)
			return
		}
		 LoginPasswordCheck(login_user, user, w);
	} else if strings.TrimSpace(login_user.Username) == ""{
        err := db.Collection_users.FindOne(context.TODO(), bson.M{
				"email": login_user.Email,
			}).Decode(&user)

		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			log.Fatal("User not found:", err)
			return
		}
		 LoginPasswordCheck(login_user, user, w);
	}
}
func LoginPasswordCheck(login_user User_Login, user models.User, w http.ResponseWriter) {
	if login_user.Password == user.Password {
                
        fmt.Println("Login Successful")

		token_user := utils.JWTUser{
		Id: user.ID.Hex(),
		UserName: user.UserName,
		Email: user.Email,
	   }

        tokens,err := utils.GenerateJWT(token_user)
	   if err != nil{
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return 
	   }
	   _ , err = db.Collection_users.UpdateOne(context.TODO(),bson.M{"_id": user.ID}, bson.M{ "$set": bson.M{"refreshToken": tokens.RefreshToken}} )
        
       if err != nil{
		http.Error(w, "Failed to Update token", http.StatusInternalServerError)
		return 
	   }

	//    data, err := json.Marshal(tokens)

	//    if err != nil {
	// 	http.Error(w, "Failed to marshal tokens", http.StatusInternalServerError)
	//    }

	//    encoded := base64.StdEncoding.EncodeToString(data)
	var buf bytes.Buffer;
    
	err = gob.NewEncoder(&buf).Encode(&tokens)
	if err != nil {
		http.Error(w, "Failed to encode cookie data", http.StatusInternalServerError)
		return 
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())


	   cookie := &http.Cookie{
		Name : "user_JWT",
		Value : encoded,
		Path: "/",
		Expires: time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure: false, // Set to true if using HTTPS
	}
       http.SetCookie(w, cookie)

	   w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(userRespone))
	   json.NewEncoder(w).Encode(map[string]interface{}{
			"_id": user.ID.Hex(),
			"fullName": user.FullName,
			"userName": user.UserName,
			"email": user.Email,
		})
		
	} else{
		http.Error(w, "Wrong Password", http.StatusUnauthorized)
		return 
 }
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(middlewares.ContextKey{})

	if userData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return 
	}
	userObjID, err := primitive.ObjectIDFromHex(userData.(string))
	fmt.Printf("%T", userObjID)
	if err!= nil{
        log.Fatal("error while coverting ID string to primitive object ID")
	}

	 _, err = db.Collection_users.UpdateOne(context.TODO(),bson.M{
		"_id": userObjID,
	 }, bson.M{
		"$unset": bson.M{"refreshToken": ""}, // Remove the refresh token field
	 })

	 if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return 
	 }

    w.Write([]byte("Logout Successfully"))
}



