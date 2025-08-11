package main

import (
	"backend-in-go/controllers"
	"backend-in-go/db"
	"backend-in-go/middlewares"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

  gob.Register(&controllers.Register_User_Cookie{})
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

  r := mux.NewRouter()

  api := r.PathPrefix("/api/v1/user").Subrouter()


  api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"));
  }) 
  
  // routes.InitUserRoutes();

  api.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
  api.HandleFunc("/login", controllers.LoginUser).Methods("POST")
  // http.HandleFunc("/filetest", cloudinary.FileDataTest)

  api.Handle("/posts", middlewares.VerifyJWT(http.HandlerFunc(controllers.Posts))).Methods("POST")

  logoutHandler := http.HandlerFunc(controllers.Logout)
  api.Handle("/logout", middlewares.VerifyJWT(logoutHandler)).Methods("POST")
  fmt.Println("Server starting on port", port)
  http.ListenAndServe(port, nil)


}