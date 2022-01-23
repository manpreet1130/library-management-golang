package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/manpreet1130/library-management/models"
	"github.com/manpreet1130/library-management/utils"
)

// GetBooks returns a list of books available in the
// database. User doesn't have to be logged in to access
// the list of books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	// User doesn't need to be logged in to view books
	books := models.GetBooks()

	res, _ := json.Marshal(books)

	log.Println("list of books sent")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// AddBook adds a new book to the collection of books
// and addition of a new book can only be done by
// user with 'admin' status
func AddBook(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		log.Println("could not find token cookie, user needs to login")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please log in to continue"))
		return
	}

	if err := models.AuthenticateUser(cookie); err != nil {
		log.Println("this user is not authorized to add a book into the system")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must be admin to add book to system"))
		return
	}

	book := &models.Book{}
	utils.ParseBody(book, r)

	addedBook := book.AddBook()

	res, _ := json.Marshal(addedBook)

	log.Printf("[ADD BOOK] Book with title %v added\n", addedBook.Title)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
