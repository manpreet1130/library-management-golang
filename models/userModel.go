package models

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/manpreet1130/library-management/database"
	"golang.org/x/crypto/bcrypt"
)

const SECRET = "secretkey"

// User consists of the following fields
// FirstName is a required field and must consist of a minimum of 2 and a maximum of 30 characters
// LastName is an optional field
// Email is a required field and must be a valid email: example@example.com
// Password is a required field and must consist of a minimum of 5 and a maximum of 100 characters
// Auth is a required field which will be filled either as 'user' or 'admin'
type User struct {
	gorm.Model
	UUID      uuid.UUID `gorm:"primary_key" json:"user_id"`
	FirstName string    `json:"firstname" validate:"required,min=2,max=30"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=5,max=100"`
	Auth      string    `json:"auth" validate:"required"`
	MyCart    Cart
}

// Check confirms that the email provided isn't already being used
// Accepts user as input and returns error if any
func (user *User) Check() error {
	db := database.GetDB()
	result := db.Where("Email = ?", user.Email).First(&User{})
	if result.RowsAffected == 0 {
		return nil
	}

	return errors.New("found user with this email address")
}

// AddUser adds the new user to the database and returns a reference
// to the newly created user
// Accepts user and returns user
func (user *User) AddUser() *User {
	db := database.GetDB()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("could not generate hashed password")
	}
	user.Password = string(hashedPassword)
	user.UUID = uuid.New()
	db.Create(&user)
	return user
}

// ValidateUser checks whether a user with the given email address is
// present in the system and further compares the given password to
// the hashed password present in the database
// Accepts user as input and returns an error if any
func (user *User) ValidateUser() error {
	db := database.GetDB()
	dbUser := &User{}
	db.Where("Email = ?", user.Email).Find(&dbUser)

	if dbUser.FirstName == "" {
		return errors.New("user with this email does not exist")
	}

	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}

// Login generates the required token string which will further be
// saved in a cookie
// takes in a reference to user and returns a string and error if any
func (user *User) Login() (string, error) {
	db := database.GetDB()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Email,
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
	})

	tokenString, err := token.SignedString([]byte(SECRET))
	dbUser := &User{}
	db.Where("Email = ?", user.Email).First(&dbUser)

	CreateCart(dbUser.UUID)

	if err != nil {
		return "", err
	}
	return tokenString, err
}

// AuthenticateUser extracts the user from the cookie generated and
// checks whether permission is granted for the user or isn't.
// Permission is only granted to the user who has authentication
// status as 'admin'
// AuthenticateUser takes in a cookie as input and returns an error if any
func AuthenticateUser(cookie *http.Cookie) error {
	db := database.GetDB()
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	user := &User{}
	db.Where("Email = ?", claims.Issuer).Find(&user)

	if user.Auth != "admin" {
		return errors.New("permission not granted for user")
	}

	return nil
}

// GetUsers returns a list of all users present in the database
func GetUsers() []User {
	db := database.GetDB()
	users := []User{}

	db.Find(&users)
	return users
}

func GetUser(cookie *http.Cookie) *User {
	db := database.GetDB()
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return nil
	}

	claims := token.Claims.(*jwt.StandardClaims)

	user := &User{}
	db.Where("Email = ?", claims.Issuer).Find(&user)

	return user
}
