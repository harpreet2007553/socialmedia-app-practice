package controllers

import (
	"backend-in-go/cloudinary"
	"backend-in-go/db"
	"backend-in-go/middlewares"
	"backend-in-go/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Posts(w http.ResponseWriter, r *http.Request) {

	userData := r.Context().Value(middlewares.ContextKey{})
	if userData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(userData.(string))
	// fmt.Printf("%T", userObjID)
	if err != nil {
		log.Fatal("error while coverting ID string to primitive object ID")
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	fmt.Println("Title:", title)

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	post := models.Post{
		Title:      title,
		Content:    content,
		Attachment: "",
		Owner:      userObjID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := db.Collection_posts.InsertOne(context.Background(), post)

	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		log.Fatal("Failed to create post:", err)
		return
	}
	postID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		log.Fatal("Failed to create post:", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		fmt.Println("error getting file:", err)
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	path := "./images/" + fileHeader.Filename
	defer file.Close()

	fmt.Println("File Name:", fileHeader.Filename)

	outFile, err := os.Create(path)
	if err != nil {
		fmt.Println("error creating file:", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}

	fmt.Println("closing file")

	// defer outFile.Close()
	fmt.Println("copying file")
	_, err = io.Copy(outFile, file)

	fmt.Println("File copied successfully")

	if err != nil {
		fmt.Println("error copying file:", err)
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}
	outFile.Close()
	cloudinary.UploadImage(path, postID, w)
	fmt.Println("hello world")
	w.Write([]byte("Post created successfully with ID: " + postID.Hex()))

}
func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	// vars:= mux.Vars(r)
	// userId := vars["userId"]
	userId := r.URL.Query().Get("userId")
	userIdObj, err := primitive.ObjectIDFromHex(userId)
    
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	pipeline := mongo.Pipeline{
	   bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userIdObj}}}},
	   bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "posts"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "owner"},
		{Key: "as", Value: "posts"},
	   }}},
	//    bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$posts"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	   bson.D{
    	{Key: "$project", Value: bson.D{
        	{Key: "fullName", Value: 1},
        	{Key: "email", Value: 1},
        	{Key: "userName", Value: 1},
        	{Key: "posts._id", Value: 1},
        	{Key: "posts.title", Value: 1},
        	{Key: "posts.content", Value: 1},
			{Key: "posts.attachment", Value: 1},
    	}},
	},

}

	cur, err := db.Collection_users.Aggregate(context.TODO(), pipeline)
    if err != nil {
        log.Fatal(err)
    }
     
	var posts []bson.M
	err = cur.All(context.TODO(), &posts)
	if err != nil{
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
    defer cur.Close(context.TODO())

	if len(posts)==0{
		w.Write([]byte("No posts found or user have no posts uploaded"))
	}

    jsonPostData, err := json.Marshal(posts)

	if err != nil {
		http.Error(w, "Failed to marshal posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPostData)
    
}

func GetCommentsOnPost(w http.ResponseWriter, r *http.Request){
	postId := r.URL.Query().Get("postId")
	if postId == ""{
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	// fmt.Println(postId)
	postIdObj, err := primitive.ObjectIDFromHex(postId)
	if err != nil{
		http.Error(w, "Error converting postId string type to primitive.ObjectId type", http.StatusBadRequest)
        return
	}
	pipeline :=  mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{ Key: "_id",Value : postIdObj}}}},
		 bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "comments"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "post_id"},
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
