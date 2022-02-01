package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/manpreet1130/library-management/models"
	"github.com/manpreet1130/library-management/utils"
)

// Signup parses the request and checks whether a user with a particular
// email already exists, and if it doesn't, a new user is created
func Signup(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	// Validation of JSON to be added.

	if err := utils.ParseBody(user, r); err != nil {
		log.Println("[SIGNUP] Unable to parse body from request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read request"))
		return
	}

	if err := user.Check(); err != nil {
		log.Println("[SIGNUP] User with this email already exists")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User with the following email already exists."))
		return
	}

	createdUser := user.AddUser()
	log.Printf("[SIGNUP] Signup successful for %v\n", createdUser.FirstName)
}

// Login parses the request containing the email and password, checks whether a user
// with the given email exists and further checks whether the password entered
// corresponds to the the hashed password saved in the database.
// Once these two conditions are confirmed, a new token is generated and saved as a cookie
func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := utils.ParseBody(user, r); err != nil {
		log.Println("[LOGIN] Unable to parse body from request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to read request sent"))
		return
	}

	if err := user.ValidateUser(); err != nil {
		log.Println("[LOGIN] User does not exist, signup required")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist, please sign up to continue"))
		return
	}

	tokenString, err := user.Login()
	if err != nil {
		log.Fatal(err)
	}

	cookie := http.Cookie{
		Name:    "Token",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 5),
	}

	http.SetCookie(w, &cookie)
	log.Println("[LOGIN] Login successful")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// GetUsers firstly checks whether there exists a cookie corresponding to the request
// sent and further authentication is done on the user as only the
// user authentication = ADMIN is allowed to receive the list of users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		log.Println("Could not find token cookie")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please log in to continue"))
		return
	}

	if err := models.AuthenticateUser(cookie); err != nil {
		log.Println("this user is not authorized to view this info")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This user does not have permission rights to view this information"))
		return
	}

	users := models.GetUsers()

	res, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Checkout performs the book checkout for a particular user
func Checkout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		log.Println("Could not find token cookie")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please log in to continue"))
		return
	}

	user := models.GetUser(cookie)
	if user == nil {
		log.Println("could not retrieve user")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while fetching user data"))
		return
	}

	user.Checkout()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Books to be returned within 7 days"))
}
