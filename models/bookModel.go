package models

import (
	"os"
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
// CartUUID: one-to-many relationship between cart and books
// Due: On checkout this field is updated
type Book struct {
	gorm.Model
	Title    string    `json:"title" validate:"required,min=2,max=50"`
	Author   string    `json:"author" validate:"required,min=2,max=50"`
	Genre    string    `json:"genre" validate:"required,min=2,max=50"`
	Quantity uint64    `json:"quantity" validate:"required,gt=0"`
	CartUUID uuid.UUID `json:"cart_uuid"`
	Due      time.Time
}

// GetBooks returns a list all books present in the database
func GetBooks() []Book {
	db := database.GetDB()
	books := []Book{}
	db.Where("cart_uuid = ?", os.Getenv("ADMIN")).Find(&books)
	return books
}

// AddBook enters a new book into the database and updates the quantity of
// book if one with the same title and author already exists
func (book *Book) AddBook() *Book {
	dbBook := &Book{}
	db := database.GetDB()

	result := db.Where("Title = ? AND Author = ? AND cart_uuid = ?", book.Title, book.Author, uuid.MustParse(os.Getenv("ADMIN"))).Find(&dbBook)

	if result.RowsAffected != 0 {
		dbBook.Quantity += book.Quantity
		db.Save(&dbBook)
		return dbBook
	}
	db.Create(&book)
	return book
}

// GetBook takes in a title as input and returns a book with the corresponding title
func GetBook(title string) *Book {
	db := database.GetDB()
	book := &Book{}
	db.Where("Title = ?", title).Find(&book)
	return book
}

// GetBooksByTitle takes in title as input and returns a list of books with the corresponding
// title present in the database
func GetBooksByTitle(title string) []Book {
	db := database.GetDB()
	books := []Book{}

	db.Where("Title = ? AND cart_uuid = ?", title, uuid.MustParse(os.Getenv("ADMIN"))).Find(&books)
	return books
}

// GetBooksByAuthor takes in author as input and returns a list of books with the corresponding
// author present in the database
func GetBooksByAuthor(author string) []Book {
	db := database.GetDB()
	books := []Book{}

	db.Where("Author = ? AND cart_uuid = ?", author, uuid.MustParse(os.Getenv("ADMIN"))).Find(&books)
	return books
}
