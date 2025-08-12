package controllers

import (
	"backend-in-go/db"
	"backend-in-go/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Like(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("postId")
	userId := r.URL.Query().Get("userId")

	postObjId, err := primitive.ObjectIDFromHex(postId)
	if err != nil{
		http.Error(w, "Post ID is required from Param",http.StatusUnauthorized)
		return
	}
    userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil{
		http.Error(w, "Post ID is required from Param",http.StatusUnauthorized)
		return
	}

	like := models.Like{
		UserID: userObjId,
        PostID: postObjId,
	}
	result ,err := db.Collection_likes.InsertOne(context.TODO(), like)

	if err != nil{
		http.Error(w, "Error while liking the post or uploading like to the database", http.StatusUnauthorized)
		return
	}
    
	likeId := result.InsertedID.(primitive.ObjectID).Hex()

	w.Write([]byte( "LikeID:" + likeId))

}
func Unlike(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("postId")
	userId := r.URL.Query().Get("userId")

	postObjId, err := primitive.ObjectIDFromHex(postId)
	if err != nil{
		http.Error(w, "Post ID is required from Param", http.StatusUnauthorized)
		return
	}
    userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil{
		http.Error(w, "Post ID is required from Param", http.StatusUnauthorized)
		return
	}

	like := models.Like{
		UserID: userObjId,
        PostID: postObjId,
	}
	_ ,err = db.Collection_likes.DeleteOne(context.TODO(), like)

	if err != nil{
		http.Error(w, "Error while liking the post or uploading like to the database", http.StatusUnauthorized)
		return
	}

	// likeId := result.DeletedCount

	w.Write([]byte( "Unliked Successfully"))

}

func GetPostLikes(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("postId")
	postObjId, err := primitive.ObjectIDFromHex(postId)

	if err != nil{
		http.Error(w, "Post ID is required from Param", http.StatusUnauthorized)
		return
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: postObjId}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "posts"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "post_id"},
			{Key: "as", Value: "likes"},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "likes", Value: 1},
			{Key: "likes_count", Value: bson.D{{Key: "$size", Value: "$likes"}}},
		}}},
	}

    cur, err := db.Collection_posts.Aggregate(context.TODO(), pipeline)
    if err != nil {
        log.Fatal(err)
    }
     
	var likes []bson.M
	err = cur.All(context.TODO(), &likes)
	if err != nil{
		http.Error(w, "Failed to retrieve likes", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(likes)
	if err != nil{
		http.Error(w, "Failed to marshal likes", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)

}
