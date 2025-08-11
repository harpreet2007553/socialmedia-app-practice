package controllers

import (
	"backend-in-go/db"
	"backend-in-go/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetUserComments(w http.ResponseWriter, r *http.Request){
   // take post id from params
   // get all comments for that post
   // return all comments as response
   userId := r.URL.Query().Get("userId")
	if userId == ""{
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	// fmt.Println(postId)
	userObjID, err := primitive.ObjectIDFromHex(userId)
	if err != nil{
		http.Error(w, "Error converting postId string type to primitive.ObjectId type", http.StatusBadRequest)
        return
	}
	pipeline :=  mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{ Key: "_id",Value : userObjID}}}},
		 bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "comments"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "owner"},
			{Key: "as", Value: "comments"},
		}}},
		// bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$comments"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "comments._id", Value: 1},
				{Key: "comments.content", Value: 1},
				{Key: "comments.owner", Value: 1},
				{Key: "comments.createdAt", Value: 1},
				{Key: "comments.updatedAt", Value: 1},
				{Key: "comments.post_id", Value: 1},
			}},
		},
	}

   cur, err := db.Collection_users.Aggregate(context.TODO(), pipeline)
    if err != nil {
        log.Fatal(err)
    }
     
	var comments []bson.M
	err = cur.All(context.TODO(), &comments)
	if err != nil{
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	if len(comments)==0{
		w.Write([]byte("No posts found or user have no posts uploaded"))
	}

    jsonCommentData, err := json.Marshal(comments)

	if err != nil {
		http.Error(w, "Failed to marshal posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCommentData)
   
}




