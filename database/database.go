package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

// Connect creates/connects to a database with the provided name
// and assigns that to the db variable
func Connect() {
	d, err := gorm.Open("sqlite3", "library.db")
	if err != nil {
		log.Fatal("could not connect to database")
	}
	db = d
}

// GetDB is used to return a reference of the database
func GetDB() *gorm.DB {
	return db
}
