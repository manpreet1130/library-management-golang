package routes

import (
	"github.com/gorilla/mux"
	"github.com/manpreet1130/library-management/controllers"
)

func CartRoutes(router *mux.Router) {
	router.HandleFunc("/users/cart", controllers.AddToCart).Methods("POST")
	router.HandleFunc("/users/cart", controllers.GetCartItems).Methods("GET")
}
