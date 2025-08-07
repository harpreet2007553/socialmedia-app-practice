package utils

import (
	// "backend-in-go/db"
	// "backend-in-go/models"
	// "context"
	"backend-in-go/db"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
)

type JWTUser struct {
	Id string
	UserName string
	Email string
}
// var founded_user models.User

type GenerateJWTResponse struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GenerateJWT(user JWTUser) (GenerateJWTResponse, error) {
	// Placeholder for JWT generation logic
	// This function should create a JWT token for the user
	
	
	// err := db.Collection_users.FindOne(context.TODO() , bson.M{"id": user.Id}).Decode(&founded_user)
	// if err!=nil {
	// 	// http.Error(nil, "user not found while verifying jwt token", http.StatusBadRequest)
	// 	log.Fatal("user not found while verifying jwt token", err)
	// 	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env variables",err)
	}
	ACCESS_TOKEN_SECRET := os.Getenv("ACCESS_TOKEN_SECRET")
	ACCESS_TOKEN_EXPIRY := os.Getenv("ACCESS_TOKEN_EXPIRY")
    AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{"_id": user.Id,"email":user.Email ,"username": user.UserName, "exp": ACCESS_TOKEN_EXPIRY})

	AccessTokenString, err := AccessToken.SignedString([]byte(ACCESS_TOKEN_SECRET))
    
	if err != nil {
		log.Fatal("error while generating access token", err)
	}
	REFRESH_TOKEN_SECRET := os.Getenv("REFRESH_TOKEN_SECRET")
	REFRESH_TOKEN_EXPIRY := os.Getenv("REFRESH_TOKEN_EXPIRY")
  
	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{"_id": user.Id, "exp": REFRESH_TOKEN_EXPIRY})

	RefreshTokenString, err := RefreshToken.SignedString([]byte(REFRESH_TOKEN_SECRET))

	if err != nil {
		log.Fatal("error while generating refresh token", err)
	}

	
	return GenerateJWTResponse{ AccessToken: AccessTokenString,RefreshToken: RefreshTokenString,} , nil

}


func VerifyJWT(AccesstokenString string, RefreshTokenString string) (string, GenerateJWTResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env variables",err)
	}
    ACCESS_TOKEN_SECRET := os.Getenv("ACCESS_TOKEN_SECRET")
	 
    _ , err = jwt.Parse(AccesstokenString, func(token *jwt.Token) (interface{}, error) {
      return []byte(ACCESS_TOKEN_SECRET), nil
   })
     
    if err != nil {
    switch {
    case errors.Is(err, jwt.ErrTokenExpired):
        fmt.Println("Token is expired")
		new_tokens, err := NewTokens(RefreshTokenString)
		if err != nil {
			log.Fatal("error while generating new tokens", err)
		}
		return "New Token Generated, Previous One Expired ",new_tokens ,nil
    case errors.Is(err, jwt.ErrTokenMalformed):
        fmt.Println("Token is malformed")
    case errors.Is(err, jwt.ErrTokenSignatureInvalid):
        fmt.Println("Invalid signature")
    default:
        fmt.Println("Other error:", err)
    }
}  

   var empty GenerateJWTResponse  
   return "Token Verified Successfully", empty ,nil
}

func NewTokens(refreshTokenString string) (GenerateJWTResponse, error) {
    err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env variables",err)
	}
	REFRESH_TOKEN_SECRET := os.Getenv("REFRESH_TOKEN_SECRET")
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(REFRESH_TOKEN_SECRET), nil
	})
	if err != nil {
		log.Fatal("error while parsing refresh token", err)
	}
	if !token.Valid {
		log.Fatal("invalid refresh token, Login Again")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Fatal("error while parsing refresh token claims")
	}
	userId, ok := claims["_id"].(string)

	userObjID , err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal("error while converting user id to object id", err)
	}

	var result_user JWTUser

	err = db.Collection_users.FindOne(context.TODO(), bson.M{ "_id": userObjID}).Decode(&result_user)
    
	if err != nil {
		log.Fatal("error while fetching user from database", err)
	}
	if !ok {
		log.Fatal("error while getting user id from refresh token claims")
	}

	NewTokens, err := GenerateJWT(result_user)
	if err != nil {
		log.Fatal("error while generating new tokens", err)
	}
	// Fetch user details from the database using userId
	return GenerateJWTResponse{AccessToken: NewTokens.AccessToken, RefreshToken: NewTokens.RefreshToken}, nil
}
