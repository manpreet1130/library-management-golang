package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/manpreet1130/library-management/models"
	"github.com/manpreet1130/library-management/utils"
)

// getUser is an internal function which is responsible for retrieving
// the cookie for token string and using that to get the user
func getUser(w http.ResponseWriter, r *http.Request) *models.User {
	cookie, err := r.Cookie("Token")
	if err != nil {
		log.Println("Could not find token cookie")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please log in to continue"))
		return nil
	}

	user := models.GetUser(cookie)
	if user == nil {
		log.Println("Could not find user")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Could not get user"))
		return nil
	}
	return user
}

// AddToCart is responsible for adding new books to the logged in user's cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	book := &models.Book{}
	if err := utils.ParseBody(book, r); err != nil {
		log.Println("could not parse book")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while parsing request"))
		return
	}

	addedBook, err := models.AddToCart(user.UUID, book)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while adding to cart"))
		return
	}

	res, _ := json.Marshal(addedBook)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetCartItems sends back all items present in the user's cart
func GetCartItems(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	books := models.GetCartItems(user.UUID)
	res, _ := json.Marshal(books)

	log.Println("[GET CART ITEMS] Sending cart items")
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// RemoveItemsFromCart removes any book which is no longer needed or  isn't to be checked out by the user
func RemoveItemsFromCart(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)

	book := &models.Book{}
	if err := utils.ParseBody(book, r); err != nil {
		log.Println("could not parse book")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while parsing request"))
		return
	}

	quant, err := models.RemoveFromCart(user.UUID, book)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This book doesn't exist in the cart"))
		return
	}

	res, _ := json.Marshal(fmt.Sprintf("%v books with title %v remain in the cart.", quant, book.Title))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
