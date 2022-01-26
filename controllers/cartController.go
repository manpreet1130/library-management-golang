package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/manpreet1130/library-management/models"
	"github.com/manpreet1130/library-management/utils"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		log.Println("Could not find token cookie")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please log in to continue"))
		return
	}

	user := models.GetUser(cookie)
	if user == nil {
		log.Println("Could not find user")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Could not get user"))
		return
	}

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
		w.Write([]byte("Error occured while adding book to cart"))
		return
	}

	res, _ := json.Marshal(addedBook)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetCartItems(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		log.Println("Could not find token cookie")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please log in to continue"))
		return
	}

	user := models.GetUser(cookie)
	if user == nil {
		log.Println("Could not find user")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Could not get user"))
		return
	}

	books := models.GetCartItems(user.UUID)

	res, _ := json.Marshal(books)

	log.Println("[GET CART ITEMS] Sending cart items")
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
