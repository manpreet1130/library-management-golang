package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/manpreet1130/library-management/database"
)

// Cart struct comprises of the following parameters
// UUID: primary key
// UserUUID: corresponds to the user who has ownership of the cart
// Books: list of books stored in the cart for later checkout
type Cart struct {
	gorm.Model
	UUID     uuid.UUID `gorm:"primary_key"`
	UserUUID uuid.UUID
	Books    []Book
}

// CreateCart takes in the user id as input and returns a reference to a newly created cart
func CreateCart(id uuid.UUID) *Cart {
	db := database.GetDB()
	cart := &Cart{
		UUID:     uuid.New(),
		UserUUID: id,
	}
	db.Save(&cart)
	return cart
}

// AddToCart takes in the user id and book to be added as input and
// returns the book added and any error
// Checks whether a cart exists for the user and creates one if it doesn't
// Further checks corresponding to the book take place and quantity of that
// book for the cart is updated if that book is already present in the cart.
// If no book with the title and author exist in the cart, a new entry is added
func AddToCart(id uuid.UUID, book *Book) (*Book, error) {
	db := database.GetDB()

	cart := &Cart{}
	result := db.Where("user_uuid = ?", id).Preload("Books").Find(&cart)
	if result.RowsAffected == 0 {
		cart = CreateCart(id)
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

// GetCartItems returns a list of all the books present in the cart corresponding
// to the user whose user id is passed as input
func GetCartItems(id uuid.UUID) []Book {
	db := database.GetDB()
	cart := &Cart{}

	db.Where("user_uuid = ?", id).Preload("Books").Find(&cart)

	return cart.Books
}

// RemoveFromCart checks for the corresponding book present in the cart and
// removes the specified quantity from the cart.
// If the quantity for any removed book comes to 0, that book is than removed from
// the cart and remaining quantity is returned.
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
