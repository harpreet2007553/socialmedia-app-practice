package controllers

import (
	"backend-in-go/db"
	"backend-in-go/models"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Follow(w http.ResponseWriter, r *http.Request){
	userId := r.URL.Query().Get("userId")
	profileId := r.URL.Query().Get("profileId")

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profileObjId, err := primitive.ObjectIDFromHex(profileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	following := models.Follow{
		Followed: profileObjId,
		FollowedBy: userObjId,
	}
    _, err = db.Collection_followings.InsertOne(context.TODO(), following)
	if err != nil {
		http.Error(w, "Error while updating follower in DB : " + err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Followed Successfully"))
}

func UnFollow(w http.ResponseWriter, r *http.Request){
	userId := r.URL.Query().Get("userId")
	profileId := r.URL.Query().Get("profileId")

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profileObjId, err := primitive.ObjectIDFromHex(profileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	unfollowing := models.Follow{
		Followed: profileObjId,
		FollowedBy: userObjId,
	}
    _, err = db.Collection_followings.DeleteOne(context.TODO(), unfollowing)
	if err != nil {
		http.Error(w, "Error while updating follower in DB : " + err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("UnFollowed Successfully"))
}

func GetUserFollowers(w http.ResponseWriter, r *http.Request){
	profileId := r.URL.Query().Get("profileId")
	profileObjId, err := primitive.ObjectIDFromHex(profileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: profileObjId},
			}},	
		},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "followings"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "followed"},
				{Key: "as", Value: "followers"},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "followers", Value: 1},
				{Key: "followers_count", Value: bson.D{
					{Key: "$size", Value: "$followers"},
				}},
			}},
		},
	}
	cur, err:= db.Collection_users.Aggregate(context.TODO(), pipeline)
	if err != nil {
		http.Error(w, "Error while getting followers in DB : " + err.Error(), http.StatusBadRequest)
		return
	}

	var followers []bson.M
    err = cur.All(context.TODO(), &followers)
	if err != nil{
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
    defer cur.Close(context.TODO())

    json.NewEncoder(w).Encode(followers)
    
}

func GetUserFollowing(w http.ResponseWriter, r *http.Request){
	profileId := r.URL.Query().Get("profileId")
	profileObjId, err := primitive.ObjectIDFromHex(profileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: profileObjId},
			}},
		},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "followings"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "followedBy"},
				{Key: "as", Value: "following"},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "following", Value: 1},
				{Key: "following_count", Value: bson.D{
					{Key: "$size", Value: "$following"},
				}},
			}},
		},
	}
	cur, err:= db.Collection_users.Aggregate(context.TODO(), pipeline)
	if err != nil {
		http.Error(w, "Error while getting following in DB : " + err.Error(), http.StatusBadRequest)
		return
	}

	var following []bson.M
    err = cur.All(context.TODO(), &following)
	if err != nil{
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
    defer cur.Close(context.TODO())

    json.NewEncoder(w).Encode(following)

}