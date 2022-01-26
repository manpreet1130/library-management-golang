package models

import (
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

	db.Where(Book{Title: book.Title, Author: book.Author}).First(&dbBook)

	// fmt.Println(cart.Books)
	cart.Books = append(cart.Books, *book)

	dbBook.Quantity -= book.Quantity

	db.Save(&cart)
	db.Save(&dbBook)

	return book, nil
}

func GetCartItems(id uuid.UUID) []Book {
	db := database.GetDB()
	cart := &Cart{}

	db.Where("user_uuid = ?", id).Preload("Books").Find(&cart)

	return cart.Books
}
