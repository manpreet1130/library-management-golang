package models

import (
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
	UserUUID uuid.UUID `json:"user_id"`
}

func GetBooks() []Book {
	db := database.GetDB()
	books := []Book{}
	db.Find(&books)
	return books
}

func (book *Book) AddBook() *Book {
	dbBook := &Book{}
	db := database.GetDB()

	result := db.Where(Book{Title: book.Title, Author: book.Author}).Find(&dbBook)

	if result.RowsAffected == 1 {
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

func UpdateBook(title string, amount uint64) {
	book := GetBook(title)
	db := database.GetDB()
	book.Quantity += amount
	db.Save(&book)
}
