package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/manpreet1130/library-management/database"
)

type Cart struct {
	gorm.Model
	UUID     uuid.UUID `gorm:"primary_key"`
	UserUUID uuid.UUID
	Books    []Book
}

func CreateCart(id uuid.UUID) {
	db := database.GetDB()
	cart := &Cart{
		UUID:     uuid.New(),
		UserUUID: id,
	}
	db.Save(&cart)
}

func AddToCart(id uuid.UUID, book *Book) (*Book, error) {
	db := database.GetDB()
	cart := &Cart{}
	dbBook := &Book{}

	db.Where("user_uuid = ?", id).Preload("Books").Find(&cart)
	if cart.UUID == uuid.MustParse(ADMIN) {
		return book, errors.New("admin logged in, log in as user")
	}

	db.Where(Book{Title: book.Title, Author: book.Author}).First(&dbBook)

	if book.Quantity > dbBook.Quantity {
		book.Quantity = dbBook.Quantity
	}

	userBook := &Book{}
	result := db.Where(Book{Title: book.Title, Author: book.Author, CartUUID: cart.UUID}).Find(&userBook)

	if result.RowsAffected != 0 {
		userBook.Quantity += book.Quantity
		db.Save(&userBook)
	} else {
		cart.Books = append(cart.Books, *book)
		db.Save(&cart)
	}

	dbBook.Quantity -= book.Quantity

	db.Save(&dbBook)

	return book, nil
}

func GetCartItems(id uuid.UUID) []Book {
	db := database.GetDB()
	cart := &Cart{}

	db.Where("user_uuid = ?", id).Preload("Books").Find(&cart)

	return cart.Books
}

// func (cart *Cart) deleteCartItem(book *Book) {
// 	db := database.GetDB()
// 	b := &Book{}
// 	db.Where(Book{Title: book.Title, Author: book.Author}).First(&b)

// 	db.Delete(&b)
// }
