package middlewares

import (
	"backend-in-go/controllers"
	"backend-in-go/db"
	"backend-in-go/utils"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTokens struct{
	AccessToken  string 
	RefreshToken string
}

func VerifyJWT(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env variables", err)
	}

	type contextKey string

    const userContextKey contextKey = "user"

	ACCESS_TOKEN_SECRET := os.Getenv("ACCESS_TOKEN_SECRET")
    
	cookie, err := r.Cookie("user_JWT")

	if err != nil {
		http.Error(w , "No JWT cookie found", http.StatusUnauthorized)
	}
    var jwt_tokens JWTokens

	data, err := base64.StdEncoding.DecodeString(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid base64", http.StatusBadRequest)
    }
	
    buf := bytes.NewBuffer(data)
    dec := gob.NewDecoder(buf)
    err = dec.Decode(&jwt_tokens)
	
    if err != nil {
        http.Error(w, "Gob decode failed", http.StatusInternalServerError)
    }


	result , err := jwt.Parse(jwt_tokens.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(ACCESS_TOKEN_SECRET), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			fmt.Println("Token is expired")
			new_tokens, err := NewTokens(w, jwt_tokens.RefreshToken)
			if err != nil {
				log.Fatal("error while generating new tokens", err)
			}
			
	 		ctx := context.WithValue(r.Context(), userContextKey, new_tokens)

			originalHandler.ServeHTTP(w, r.WithContext(ctx))
			return
			
			// return "New Token Generated, Previous One Expired ", new_tokens, nil
		case errors.Is(err, jwt.ErrTokenMalformed):
			http.Error(w, "Malformed token", http.StatusBadRequest)
			return
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):

			http.Error(w, "Invalid token signature", http.StatusUnauthorized)
			return
		default:
			http.Error(w, "Error parsing token", http.StatusInternalServerError)
			return
		}
	}
	if !result.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

    
	jwt_map , ok := result.Claims.(jwt.MapClaims)

	if !ok {
		log.Fatal("error while parsing jwt claims")
	}
	// userId, ok := jwt_map["_id"].(string)

	
     
	ctx := context.WithValue(r.Context(), userContextKey, jwt_map)

	originalHandler.ServeHTTP(w, r.WithContext(ctx))
    
    // var empty utils.GenerateJWTResponse
     
	// next.ServeHTTP(w, r.WithContext(ctx))
	// return "Token Verified Successfully", empty, nil
})
}

func NewTokens(w http.ResponseWriter,refreshTokenString string) (utils.GenerateJWTResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env variables", err)
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

	userObjID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal("error while converting user id to object id", err)
	}

	var result_user utils.JWTUser

	err = db.Collection_users.FindOne(context.TODO(), bson.M{"_id": userObjID}).Decode(&result_user)

	if err != nil {
		log.Fatal("error while fetching user from database", err)
	}
	if !ok {
		log.Fatal("error while getting user id from refresh token claims")
	}

	NewTokens, err := utils.GenerateJWT(result_user)
	if err != nil {
		log.Fatal("error while generating new tokens", err)
	}

	_ , err = db.Collection_users.UpdateOne(context.TODO(),bson.M{"_id": userObjID}, bson.M{ "$set": bson.M{"refreshToken": NewTokens.RefreshToken}} )
    
	if err != nil {
		http.Error(w, "Failed to Update token", http.StatusInternalServerError)
		return utils.GenerateJWTResponse{}, err
	}

	var cookie_data = controllers.Register_User_Cookie{
		RefreshToken:     NewTokens.RefreshToken,
		AccessToken:     NewTokens.AccessToken,
	}
	var buf bytes.Buffer;
    
	err = gob.NewEncoder(&buf).Encode(&cookie_data)
	if err != nil {
		http.Error(w, "Failed to encode cookie data", http.StatusInternalServerError)
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

	// Fetch user details from the database using userId
	return utils.GenerateJWTResponse{AccessToken: NewTokens.AccessToken, RefreshToken: NewTokens.RefreshToken}, nil
}