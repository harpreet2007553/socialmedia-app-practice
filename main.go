package main

import (
	"fmt"
	"net/http"
	"os"
  "BackendWithGolang/db"
	"github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load();
  if err!= nil {
    fmt.Println("Error loading .env file")
  }

  port := os.Getenv("PORT") ;

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
  })

  http.ListenAndServe(port, nil)
  fmt.Println("Hello World!")

  ConnectDB()

}