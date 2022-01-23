package routes

import (
	"github.com/gorilla/mux"
	"github.com/manpreet1130/library-management/controllers"
)

func UserRoutes(router *mux.Router) {
	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
}
