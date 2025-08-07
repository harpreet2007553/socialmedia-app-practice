package utils

import (
	// "backend-in-go/db"
	// "backend-in-go/models"
	// "context"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

	

func VerifyJWT(AccesstokenString string) (string, error) {
	 err := godotenv.Load()
	 if err != nil {
		log.Fatal("error loading env variables",err)
	}
	 ACCESS_TOKEN_SECRET := os.Getenv("ACCESS_TOKEN_SECRET")
	 
      token, err := jwt.Parse(AccesstokenString, func(token *jwt.Token) (interface{}, error) {
      return ACCESS_TOKEN_SECRET, nil
   })
  
   if err != nil {
      return "Invalid Token", err
   }
  
   if !token.Valid {
      return "Invalid Token", nil
   }
  
   return "Verified Successfully",nil
}