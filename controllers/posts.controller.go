package controllers

import (
	"backend-in-go/cloudinary"
	"backend-in-go/db"
	"backend-in-go/middlewares"
	"backend-in-go/models"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
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