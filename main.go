package main

import (
	"backend-in-go/controllers"
	"backend-in-go/db"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
  db.ConnectDB();
  err := godotenv.Load();
  if err != nil {
    fmt.Println("Error loading .env file")
  }

  port := os.Getenv("PORT");
  if port == "" {
    port = "8080"
  }
  port = ":" + port
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"));
  }) 
  
  // routes.InitUserRoutes();

  http.HandleFunc("/register", controllers.RegisterUser)

  fmt.Println("Server starting on port", port)
  http.ListenAndServe(port, nil)


}