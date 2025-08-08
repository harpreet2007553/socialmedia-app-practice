package main

import (
	"backend-in-go/controllers"
	"backend-in-go/db"
	"backend-in-go/middlewares"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

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
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"));
  }) 
  
  // routes.InitUserRoutes();

  http.HandleFunc("/register", controllers.RegisterUser)
  http.HandleFunc("/login", controllers.LoginUser)

  logoutHandler := http.HandlerFunc(controllers.Logout)
  http.Handle("/logout", middlewares.VerifyJWT(logoutHandler))
  fmt.Println("Server starting on port", port)
  http.ListenAndServe(port, nil)


}