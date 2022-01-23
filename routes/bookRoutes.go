package routes

import (
	"github.com/gorilla/mux"
	"github.com/manpreet1130/library-management/controllers"
)

func BookRoutes(router *mux.Router) {
	router.HandleFunc("/books", controllers.GetBooks).Methods("GET")
	router.HandleFunc("/books", controllers.AddBook).Methods("POST")
	// router.HandleFunc("/books", controllers.GetBooksByTitle).Methods("GET")
	// router.HandleFunc("/books/{author}", controllers.GetBooksByAuthor).Methods("GET")
}
