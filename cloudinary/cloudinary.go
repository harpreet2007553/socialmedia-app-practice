package cloudinary

import (
	"backend-in-go/db"
	"context"
	"encoding/json"
	"fmt"
	"log"

	// "io"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func UploadImage(ImageLocalPath string, PostId primitive.ObjectID, w http.ResponseWriter) {

	// defer removeFile(ImageLocalPath)
	PostIdStr := PostId.Hex()
	// fmt.Println("Bug Check")
	err:= godotenv.Load()
	if err != nil {
		http.Error(w, "Error loading environment variables", http.StatusInternalServerError)
		log.Fatal("error loading env variables", err)
	}

	cldUrl := os.Getenv("CLOUDINARY_URL")
	fmt.Println("CLOUDINARY_URL:", cldUrl)
    cld, _ := cloudinary.NewFromURL(cldUrl)

    cld.Config.URL.Secure = true
    ctx := context.Background()
  // Upload the image.
  // Set the asset's public ID and allow overwriting the asset with new versions
    resp, err := cld.Upload.Upload(ctx, ImageLocalPath, uploader.UploadParams{
        PublicID:       "post_" + PostIdStr,
        UniqueFilename: api.Bool(true),
        Overwrite:      api.Bool(true)})
    if err != nil {
        fmt.Println(err)
    }
	
	imgUrl := resp.SecureURL

	pubId := resp.PublicID

	attachment := map[string]string{
		"public_id": pubId,
		"secure_url": imgUrl,
	}

	jsonAttachment, err := json.Marshal(attachment)
	if err != nil {
		fmt.Println("error marshalling attachment:", err)
		return
	}

	_, err = db.Collection_posts.UpdateOne(context.TODO(), bson.M{
		"_id": PostId,
	},
    bson.M{
	"$set": bson.M{
		"attachment": imgUrl,
	    },
   })

   	if err != nil {
		http.Error(w, "Failed to update post with image URL", http.StatusInternalServerError)
		log.Fatal("Failed to update post with image URL:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonAttachment)

  // Log the delivery URL
    fmt.Println("****2. Upload an image****\nDelivery URL:", resp.SecureURL)

	err = os.Remove(ImageLocalPath)
	if err != nil {
		fmt.Println("error deleting file:", err)
		
	}
	
	
}
