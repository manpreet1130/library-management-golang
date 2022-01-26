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

func CreateCart(id uuid.UUID) *Cart {
	db := database.GetDB()
	cart := &Cart{
		UUID:     uuid.New(),
		UserUUID: id,
	}
	db.Save(&cart)
	return cart
}

func AddToCart(id uuid.UUID, book *Book) (*Book, error) {
	db := database.GetDB()

	cart := &Cart{}
	result := db.Where("user_uuid = ?", id).Preload("Books").Find(&cart)
	if result.RowsAffected == 0 {
		cart = CreateCart(id)
	}

	if cart.UUID == uuid.MustParse(ADMIN) {
		return book, errors.New("admin logged in, log in as user")
	}

	dbBook := &Book{}
	db.Where("Title = ? AND Author = ?", book.Title, book.Author).First(&dbBook)

	if dbBook.Quantity == 0 {
		return book, errors.New("book out of stock")
	}

	if book.Quantity > dbBook.Quantity {
		book.Quantity = dbBook.Quantity
	}

	userBook := &Book{}
	result = db.Where("Title = ? AND Author = ? AND cart_uuid = ?", book.Title, book.Author, cart.UUID).Find(&userBook)

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

func RemoveFromCart(id uuid.UUID, book *Book) (uint64, error) {
	db := database.GetDB()

	cart := &Cart{}
	db.Where("user_uuid = ?", id).Find(&cart)

	dbBook := &Book{}
	db.Where("Title = ? AND Author = ?", book.Title, book.Author).First(&dbBook)

	cartBook := &Book{}
	result := db.Where("Title = ? AND Author = ? AND cart_uuid = ?", book.Title, book.Author, cart.UUID).Find(&cartBook)

	if result.RowsAffected == 0 {
		return 0, errors.New("this book does not exist in the cart")
	}

	if book.Quantity > cartBook.Quantity {
		dbBook.Quantity += cartBook.Quantity
		cartBook.Quantity = 0
	} else {
		cartBook.Quantity -= book.Quantity
		dbBook.Quantity += book.Quantity
	}

	if cartBook.Quantity == 0 {
		db.Delete(&cartBook)
		return 0, nil
	}

	db.Save(&cartBook)
	db.Save(&dbBook)
	return cartBook.Quantity, nil

}

func EmptyCart(id uuid.UUID) {
	// db := database.GetDB()
	// cart := &Cart{}
	// db.Where("user_uuid = ?", id).Find(&cart)
	// db.Delete(&cart)
}
