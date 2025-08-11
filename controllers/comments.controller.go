package controllers

import (
	"backend-in-go/db"
	"backend-in-go/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Comment(w http.ResponseWriter, r *http.Request){
   // take post id from params
   // create a comment struct object
   // save comment to db
   // return the comment details as response

   // vars := mux.Vars(r)
   
   userId := r.URL.Query().Get("userId")
   postId := r.URL.Query().Get("postId")
   
   fmt.Println(userId)

   userObjID, err := primitive.ObjectIDFromHex(userId)
   if err != nil {
	   http.Error(w, "Error while converting UserID to primitive.ObjectID type", http.StatusBadRequest)
	   return
   }
   postObjID, err := primitive.ObjectIDFromHex(postId)
   if err != nil {
	   http.Error(w, "Error while converting PostID to primitive.ObjectID type", http.StatusBadRequest)
	   return
   }


   comment := models.Comment{
	   Content : r.FormValue("content"),
	   Owner : userObjID,
	   PostID : postObjID,
	   CreatedAt: time.Now(),
	   UpdatedAt: time.Now(),
   }
   
   result, err := db.Collection_comments.InsertOne(context.TODO(), comment)
   if err != nil {
	   http.Error(w, "Failed to create comment", http.StatusInternalServerError)
	   return
   }
   InsertedID := result.InsertedID.(primitive.ObjectID)

   commentResp := map[string]interface{}{
	   "CommentID": InsertedID.Hex(),
	   "Result": "Commented Successfully",
   }
   jsonComment , err := json.Marshal(commentResp)
   if err != nil {
	   http.Error(w, "Failed to marshal comment response", http.StatusInternalServerError)
	   return
   }

   w.Header().Set("Content-Type", "application/json")
   w.Write(jsonComment)
}



