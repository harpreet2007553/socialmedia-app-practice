package routes

import (
	"backend-in-go/controllers"

	"github.com/gorilla/mux"
)

func InitUserRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.RegisterUser).Methods("GET")
    
}