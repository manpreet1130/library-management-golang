package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/manpreet1130/library-management/database"
)

// Book structure consists of the following fields
// Title: required field, with a minimum and maximum of 2 and 50 characters respectively
// Author: required field
// Genre: required field
// Quantity: required field
type Book struct {
	gorm.Model
	Title    string    `json:"title" validate:"required,min=2,max=50"`
	Author   string    `json:"author" validate:"required,min=2,max=50"`
	Genre    string    `json:"genre" validate:"required,min=2,max=50"`
	Quantity uint64    `json:"quantity" validate:"required,gt=0"`
	CartUUID uuid.UUID `json:"cart_uuid"`
	Due      time.Time
}

const ADMIN = "00000000-0000-0000-0000-000000000000"

func GetBooks() []Book {
	db := database.GetDB()
	books := []Book{}
	db.Where("cart_uuid = ?", ADMIN).Find(&books)
	return books
}

func (book *Book) AddBook() *Book {
	dbBook := &Book{}
	db := database.GetDB()

	result := db.Where("Title = ? AND Author = ? AND cart_uuid = ?", book.Title, book.Author, uuid.MustParse(ADMIN)).Find(&dbBook)

	// fmt.Println(result.RowsAffected)

	if result.RowsAffected != 0 {
		dbBook.Quantity += book.Quantity
		db.Save(&dbBook)
		return dbBook
	}
	db.Create(&book)
	return book
}

func GetBook(title string) *Book {
	db := database.GetDB()
	book := &Book{}
	db.Where("Title = ?", title).Find(&book)
	return book
}

func GetBooksByTitle(title string) []Book {
	db := database.GetDB()
	books := []Book{}

	db.Where("Title = ? AND cart_uuid = ?", title, uuid.MustParse(ADMIN)).Find(&books)
	return books
}

func GetBooksByAuthor(author string) []Book {
	db := database.GetDB()
	books := []Book{}

	db.Where("Author = ? AND cart_uuid = ?", author, uuid.MustParse(ADMIN)).Find(&books)
	return books
}
